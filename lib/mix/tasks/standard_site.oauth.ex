defmodule Mix.Tasks.StandardSite.Oauth do
  use Mix.Task

  @moduledoc """
  Authorize Standard.site publishing with ATProto OAuth.

      mix standard_site.oauth --handle pdewey.com

  The task requests `repo:site.standard.document?action=create&action=update`,
  opens a browser, listens for the localhost callback, stores OAuth tokens in a
  local DETS store, and writes the active session key under
  `priv/standard_site_oauth/`.
  """

  @shortdoc "Authorize this site publisher with ATProto OAuth"

  @requirements ["app.start"]
  @callback_timeout 300_000

  @impl Mix.Task
  def run(args) do
    {opts, _args, _invalid} =
      OptionParser.parse(args,
        strict: [
          handle: :string,
          port: :integer,
          scope: :keep,
          no_open: :boolean
        ]
      )

    handle = opts[:handle] || System.get_env("ATPROTO_IDENTIFIER") || "pdewey.com"
    port = opts[:port] || Site.StandardSite.OAuth.default_port()
    scopes = Keyword.get_values(opts, :scope)

    Site.StandardSite.OAuth.configure!(port: port, scopes: scopes)
    Mix.shell().info("Requesting OAuth scopes: #{Atex.Config.OAuth.scopes()}")

    {:ok, listener} = listen(port)
    Mix.shell().info("Listening for OAuth callback on http://127.0.0.1:#{port}/oauth/callback")

    callback = Task.async(fn -> accept_callback(listener) end)

    {authz_metadata, code_verifier, authz_url} = authorization_url!(handle)

    Mix.shell().info("Open this URL to authorize #{handle}:\n\n#{authz_url}\n")
    maybe_open_browser(authz_url, opts[:no_open])

    params = Task.await(callback, @callback_timeout + 5_000)
    exchange_callback!(params, authz_metadata, code_verifier)
  end

  defp listen(port) do
    :gen_tcp.listen(port, [
      :binary,
      packet: :raw,
      active: false,
      reuseaddr: true,
      ip: {127, 0, 0, 1}
    ])
  end

  defp accept_callback(listener) do
    {:ok, socket} = :gen_tcp.accept(listener, @callback_timeout)
    {:ok, request} = :gen_tcp.recv(socket, 0, 5_000)
    params = parse_callback_params(request)

    body = "OAuth complete. You can close this tab and return to the terminal.\n"

    :ok =
      :gen_tcp.send(socket, [
        "HTTP/1.1 200 OK\r\n",
        "content-type: text/plain\r\n",
        "content-length: ",
        Integer.to_string(byte_size(body)),
        "\r\nconnection: close\r\n\r\n",
        body
      ])

    :gen_tcp.close(socket)
    :gen_tcp.close(listener)
    params
  end

  defp parse_callback_params(request) do
    [request_line | _headers] = String.split(request, "\r\n", parts: 2)
    ["GET", target | _] = String.split(request_line, " ")
    uri = URI.parse(target)

    if uri.path != "/oauth/callback" do
      Mix.raise("Unexpected OAuth callback path #{inspect(uri.path)}")
    end

    URI.decode_query(uri.query || "")
  end

  defp authorization_url!(handle) do
    with {:ok, identity} <- Atex.IdentityResolver.resolve(handle),
         pds when is_binary(pds) <- Atex.DID.Document.get_pds_endpoint(identity.document),
         {:ok, authz_server} <- Atex.OAuth.Discovery.get_authorization_server(pds),
         {:ok, authz_metadata} <-
           Atex.OAuth.Discovery.get_authorization_server_metadata(authz_server) do
      state = Atex.OAuth.create_nonce()
      code_verifier = Atex.OAuth.create_nonce()
      Process.put(:oauth_state, state)
      Process.put(:oauth_pds, pds)

      case Atex.OAuth.Flow.create_authorization_url(
             authz_metadata,
             state,
             code_verifier,
             handle
           ) do
        {:ok, url} ->
          {authz_metadata, code_verifier, url}

        {:error, reason} ->
          Mix.raise("Failed to create OAuth authorization URL: #{inspect(reason)}")
      end
    else
      error -> Mix.raise("Failed to resolve OAuth metadata for #{handle}: #{inspect(error)}")
    end
  end

  defp maybe_open_browser(_url, true), do: :ok

  defp maybe_open_browser(url, _no_open?) do
    case System.find_executable("xdg-open") do
      nil -> :ok
      xdg_open -> System.cmd(xdg_open, [url], stderr_to_stdout: true)
    end

    :ok
  end

  defp exchange_callback!(params, authz_metadata, code_verifier) do
    stored_state = Process.get(:oauth_state)
    pds = Process.get(:oauth_pds)

    unless params["state"] == stored_state and is_binary(params["code"]) do
      Mix.raise("OAuth callback was missing code or had mismatched state")
    end

    dpop_key = JOSE.JWK.generate_key({:ec, "P-256"})

    with {:ok, tokens, dpop_nonce} <-
           Atex.OAuth.Flow.validate_authorization_code(
             authz_metadata,
             dpop_key,
             params["code"],
             code_verifier
           ),
         {:ok, identity} <- Atex.IdentityResolver.resolve(tokens.did),
         ^pds <- Atex.DID.Document.get_pds_endpoint(identity.document) do
      session = %Atex.OAuth.Session{
        iss: authz_metadata.issuer,
        aud: pds,
        sub: tokens.did,
        nonce: Atex.OAuth.create_nonce(),
        access_token: tokens.access_token,
        refresh_token: tokens.refresh_token,
        expires_at: tokens.expires_at,
        dpop_key: dpop_key,
        dpop_nonce: dpop_nonce
      }

      session_key = Atex.OAuth.SessionStore.session_key(session)
      :ok = Atex.OAuth.SessionStore.insert(session)
      :dets.sync(:atex_oauth_sessions)
      Site.StandardSite.OAuth.write_session_key!(session_key)

      Mix.shell().info("OAuth session stored for #{tokens.did}")
      Mix.shell().info("Session key written to #{Site.StandardSite.OAuth.session_key_path()}")
    else
      error -> Mix.raise("Failed to exchange OAuth callback: #{inspect(error)}")
    end
  end
end
