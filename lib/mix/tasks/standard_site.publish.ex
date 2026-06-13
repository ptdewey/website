defmodule Mix.Tasks.StandardSite.Publish do
  use Mix.Task

  @moduledoc """
  Publish local blog posts as Standard.site document records.

      mix standard_site.publish --dry-run
      mix standard_site.publish --auth oauth
      mix standard_site.publish --auth app-password

  OAuth sessions are created with `mix standard_site.oauth`. App-password auth
  reads `ATPROTO_IDENTIFIER` and `ATPROTO_APP_PASSWORD`.
  """

  @shortdoc "Publish local blog posts as Standard.site document records"

  alias Site.StandardSite
  alias Site.StandardSite.Documents

  @requirements ["app.start"]

  @impl Mix.Task
  def run(args) do
    {opts, args, _invalid} =
      OptionParser.parse(args,
        strict: [
          did: :string,
          publication: :string,
          pds: :string,
          identifier: :string,
          password: :string,
          output: :string,
          post: :string,
          auth: :string,
          oauth_session_key: :string,
          oauth_port: :integer,
          dry_run: :boolean,
          validate: :boolean
        ],
        aliases: [n: :dry_run]
      )

    opts = Map.new(opts)
    opts = if "--dry-run" in args, do: Map.put(opts, :dry_run, true), else: opts

    did = opts[:did] || System.get_env("STANDARD_SITE_DID") || StandardSite.did()

    publication =
      opts[:publication] || System.get_env("STANDARD_SITE_PUBLICATION") ||
        StandardSite.publication()

    pds = opts[:pds] || System.get_env("STANDARD_SITE_PDS") || Site.AtProto.resolve_pds!(did)
    output = opts[:output] || StandardSite.documents_path()
    dry_run? = Map.get(opts, :dry_run, false)
    validate? = Map.get(opts, :validate, false)

    posts = posts(opts[:post])
    local_documents = Documents.all(output)

    remote_records =
      did
      |> Mix.Tasks.StandardSite.Sync.fetch_documents(publication, pds)

    localized_remote_records = Mix.Tasks.StandardSite.Sync.localize_record_paths(remote_records)

    remote_documents =
      localized_remote_records
      |> Documents.merge_records(%{})

    remote_records_by_path =
      remote_records
      |> Enum.map(fn %{"value" => value} ->
        {Mix.Tasks.StandardSite.Sync.local_path_for_record_value(value), value}
      end)
      |> Enum.filter(fn {path, _value} -> is_binary(path) end)
      |> Map.new()

    documents = Map.merge(local_documents, remote_documents)

    plans = Enum.map(posts, &plan_post(&1, documents, remote_records_by_path, publication))
    print_plan(plans, dry_run?)

    if dry_run? or Enum.all?(plans, &(&1.action == :unchanged)) do
      :ok
    else
      client = auth_client!(did, pds, opts)

      {documents, _client} =
        Enum.reduce(plans, {documents, client}, fn plan, {documents, client} ->
          case plan.action do
            :unchanged ->
              {documents, client}

            action when action in [:create, :update] ->
              record = prepare_record_for_publish(plan)
              {response, client} = put_record!(client, did, plan.rkey, record, validate?)

              attrs = %{
                uri: response.body["uri"],
                cid: response.body["cid"],
                rkey: plan.rkey,
                title: plan.post.title,
                published_at: record["publishedAt"],
                updated_at: record["updatedAt"]
              }

              Mix.shell().info("#{action}d #{plan.post.path} -> #{attrs.uri}")
              {Map.put(documents, plan.post.path, attrs), client}
          end
        end)

      Documents.write!(documents, output)
      Mix.shell().info("Updated #{output}")
    end
  end

  defp posts(nil), do: Site.Blog.all_posts()

  defp posts(post_id_or_path) do
    Site.Blog.all_posts()
    |> Enum.filter(fn post -> post.id == post_id_or_path or post.path == post_id_or_path end)
    |> case do
      [] -> Mix.raise("No post matched #{inspect(post_id_or_path)}")
      posts -> posts
    end
  end

  defp plan_post(post, documents, remote_records_by_path, publication) do
    metadata = Map.get(documents, post.path, %{})
    remote_record = Map.get(remote_records_by_path, post.path)
    record = StandardSite.record_for_post(post, publication: publication)
    current_record = comparable_record(remote_record)
    next_record = comparable_record(record)

    rkey =
      metadata[:rkey] || StandardSite.rkey_from_uri(metadata[:uri]) ||
        StandardSite.slug_rkey(post)

    action =
      cond do
        is_nil(remote_record) and metadata == %{} -> :create
        is_nil(remote_record) -> :update
        current_record == next_record -> :unchanged
        true -> :update
      end

    %{
      action: action,
      post: post,
      rkey: rkey,
      record: record,
      remote_record: remote_record,
      current_record: current_record,
      next_record: next_record,
      metadata: metadata
    }
  end

  defp comparable_record(nil), do: %{}

  defp comparable_record(record) when is_map(record) do
    Map.take(record, [
      "site",
      "path",
      "title",
      "description",
      "publishedAt",
      "content",
      "textContent",
      "tags",
      "bskyPostRef"
    ])
  end

  defp prepare_record_for_publish(%{action: :update, record: record}) do
    Map.put(
      record,
      "updatedAt",
      DateTime.utc_now() |> DateTime.truncate(:second) |> DateTime.to_iso8601()
    )
  end

  defp prepare_record_for_publish(%{record: record}), do: record

  defp print_plan(plans, dry_run?) do
    suffix = if dry_run?, do: " (dry-run)", else: ""
    Mix.shell().info("Standard.site publish plan#{suffix}:")

    Enum.each(plans, fn plan ->
      Mix.shell().info("  #{plan.action} #{plan.post.path} rkey=#{plan.rkey}")

      if dry_run? and plan.action in [:create, :update] do
        print_record_diff(plan)
      end
    end)
  end

  defp print_record_diff(%{action: :create, next_record: next}) do
    next
    |> Enum.sort_by(fn {key, _value} -> key end)
    |> Enum.each(fn {key, value} ->
      Mix.shell().info("    + #{key}: #{summarize_value(value)}")
    end)
  end

  defp print_record_diff(%{current_record: current, next_record: next}) do
    (Map.keys(current) ++ Map.keys(next))
    |> Enum.uniq()
    |> Enum.sort()
    |> Enum.each(fn key ->
      old = Map.get(current, key)
      new = Map.get(next, key)

      if old != new do
        Mix.shell().info("    ~ #{key}:")
        Mix.shell().info("      - #{summarize_value(old)}")
        Mix.shell().info("      + #{summarize_value(new)}")
      end
    end)
  end

  defp summarize_value(nil), do: "<missing>"

  defp summarize_value(%{"$type" => type} = value) do
    details =
      value
      |> Map.delete("$type")
      |> Enum.sort_by(fn {key, _value} -> key end)
      |> Enum.map_join(", ", fn {key, nested_value} ->
        "#{key}=#{summarize_nested_value(nested_value)}"
      end)

    "%{\"$type\" => #{inspect(type)}, #{details}}"
  end

  defp summarize_value(value) when is_map(value) do
    value
    |> Enum.sort_by(fn {key, _value} -> key end)
    |> Enum.map_join(", ", fn {key, nested_value} ->
      "#{key}=#{summarize_nested_value(nested_value)}"
    end)
    |> then(&"%{#{&1}}")
  end

  defp summarize_value(value) when is_list(value), do: inspect(value)

  defp summarize_value(value) when is_binary(value) do
    value
    |> String.replace(~r/\s+/, " ")
    |> String.trim()
    |> truncate(160)
    |> inspect()
  end

  defp summarize_value(value), do: inspect(value)

  defp summarize_nested_value(value) when is_binary(value) do
    value
    |> String.replace(~r/\s+/, " ")
    |> String.trim()
    |> then(fn text -> "#{byte_size(value)} bytes #{inspect(truncate(text, 80))}" end)
  end

  defp summarize_nested_value(value), do: summarize_value(value)

  defp truncate(text, max_length) when byte_size(text) <= max_length, do: text

  defp truncate(text, max_length) do
    String.slice(text, 0, max_length) <> "…"
  end

  defp auth_client!(did, pds, opts) do
    auth = opts[:auth] || System.get_env("STANDARD_SITE_AUTH") || default_auth(opts)

    case auth do
      "oauth" -> oauth_client!(did, opts)
      "app-password" -> login_client!(pds, opts)
      other -> Mix.raise("Unknown auth mode #{inspect(other)}; use oauth or app-password")
    end
  end

  defp default_auth(opts) do
    cond do
      opts[:oauth_session_key] || System.get_env("STANDARD_SITE_OAUTH_SESSION_KEY") -> "oauth"
      match?({:ok, _}, Site.StandardSite.OAuth.session_key()) -> "oauth"
      true -> "app-password"
    end
  end

  defp oauth_client!(did, opts) do
    Site.StandardSite.OAuth.configure!(
      port: opts[:oauth_port] || Site.StandardSite.OAuth.default_port()
    )

    session_key =
      opts[:oauth_session_key] ||
        case Site.StandardSite.OAuth.session_key() do
          {:ok, session_key} ->
            session_key

          {:error, :not_found} ->
            Mix.raise("No OAuth session found; run mix standard_site.oauth first")
        end

    case Atex.XRPC.OAuthClient.new(session_key) do
      {:ok, client} ->
        oauth_did = Atex.XRPC.OAuthClient.did(client)

        if oauth_did != did do
          Mix.raise("OAuth session DID #{oauth_did} does not match publish DID #{did}")
        end

        client

      {:error, reason} ->
        Mix.raise("""
        Could not load OAuth session #{inspect(session_key)}: #{inspect(reason)}

        The saved session key exists, but the local OAuth token store does not
        have a matching session. Re-authorize with:

            mix standard_site.oauth --handle pdewey.com

        If that keeps happening, remove the stale local files and try again:

            rm -rf priv/standard_site_oauth priv/dets
            mix standard_site.oauth --handle pdewey.com
        """)
    end
  end

  defp login_client!(pds, opts) do
    identifier =
      opts[:identifier] || System.get_env("ATPROTO_IDENTIFIER") ||
        System.get_env("STANDARD_SITE_IDENTIFIER")

    password =
      opts[:password] || System.get_env("ATPROTO_APP_PASSWORD") ||
        System.get_env("STANDARD_SITE_APP_PASSWORD")

    unless identifier && password do
      Mix.raise(
        "Set ATPROTO_IDENTIFIER and ATPROTO_APP_PASSWORD, or pass --identifier and --password"
      )
    end

    case Atex.XRPC.LoginClient.login(pds, identifier, password) do
      {:ok, client} -> client
      {:error, error} -> Mix.raise("Failed to log in to #{pds}: #{inspect(error)}")
      other -> Mix.raise("Failed to log in to #{pds}: #{inspect(other)}")
    end
  end

  defp put_record!(client, did, rkey, record, validate?) do
    json = %{
      repo: did,
      collection: Site.StandardSite.collection(),
      rkey: rkey,
      record: record,
      validate: validate?
    }

    case Atex.XRPC.post(client, "com.atproto.repo.putRecord", json: json) do
      {:ok, response, client} -> {response, client}
      {:error, error, _client} -> Mix.raise("Failed to put #{rkey}: #{inspect(error)}")
      {:error, error} -> Mix.raise("Failed to put #{rkey}: #{inspect(error)}")
    end
  end
end
