defmodule Site.StandardSite.OAuth do
  @moduledoc false

  @dir "priv/standard_site_oauth"
  @key_path Path.join(@dir, "private_key.der.b64")
  @session_key_path Path.join(@dir, "session_key")
  @key_id "pdewey-website-cli"
  @default_port 8765
  @default_scopes ["repo:site.standard.document?action=create&action=update"]

  def dir, do: @dir
  def session_key_path, do: @session_key_path
  def default_port, do: @default_port
  def default_scopes, do: @default_scopes

  def configure!(opts \\ []) do
    port = Keyword.get(opts, :port, @default_port)
    scopes = opts |> Keyword.get(:scopes, []) |> normalize_scopes()

    Application.put_env(:atex, Atex.OAuth,
      base_url: "http://127.0.0.1:#{port}/oauth",
      is_localhost: true,
      scopes: scopes,
      private_key: private_key_b64!(),
      key_id: @key_id
    )
  end

  defp normalize_scopes(scopes) do
    (@default_scopes ++ List.wrap(scopes))
    |> Enum.reject(&(&1 in [nil, ""]))
    |> Enum.uniq()
  end

  def session_key do
    cond do
      key = System.get_env("STANDARD_SITE_OAUTH_SESSION_KEY") ->
        {:ok, String.trim(key)}

      File.exists?(@session_key_path) ->
        {:ok, @session_key_path |> File.read!() |> String.trim()}

      true ->
        {:error, :not_found}
    end
  end

  def write_session_key!(session_key) do
    File.mkdir_p!(@dir)
    File.write!(@session_key_path, session_key <> "\n")
  end

  defp private_key_b64! do
    cond do
      key = System.get_env("STANDARD_SITE_OAUTH_PRIVATE_KEY") ->
        key

      File.exists?(@key_path) ->
        @key_path |> File.read!() |> String.trim()

      true ->
        key = new_private_key_b64()
        File.mkdir_p!(@dir)
        File.write!(@key_path, key <> "\n")
        key
    end
  end

  defp new_private_key_b64 do
    {_kty, der} =
      {:ec, "P-256"}
      |> JOSE.JWK.generate_key()
      |> JOSE.JWK.to_der()

    Base.encode64(der)
  end
end
