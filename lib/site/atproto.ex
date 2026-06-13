defmodule Site.AtProto do
  @moduledoc false

  @plc_directory "https://plc.directory"
  @collection "site.standard.document"

  def list_site_documents(opts) do
    opts
    |> Keyword.put(:collection, @collection)
    |> list_records()
  end

  def list_records(opts) do
    did = Keyword.fetch!(opts, :did)
    collection = Keyword.fetch!(opts, :collection)
    pds = Keyword.get(opts, :pds) || resolve_pds!(did)

    list_records(pds, did, collection)
  end

  def resolve_pds!("did:plc:" <> _ = did) do
    did
    |> then(&get_json!("#{@plc_directory}/#{URI.encode(&1)}"))
    |> pds_from_did_document!()
  end

  def resolve_pds!("did:web:" <> domain) do
    domain = String.replace(domain, ":", "/")

    "https://#{domain}/.well-known/did.json"
    |> get_json!()
    |> pds_from_did_document!()
  end

  defp list_records(pds, did, collection, cursor \\ nil, records \\ []) do
    query =
      [repo: did, collection: collection, limit: 100, cursor: cursor]
      |> Enum.reject(fn {_key, value} -> is_nil(value) end)
      |> URI.encode_query()

    response =
      get_json!("#{String.trim_trailing(pds, "/")}/xrpc/com.atproto.repo.listRecords?#{query}")

    records = records ++ Map.fetch!(response, "records")

    case response["cursor"] do
      nil -> records
      "" -> records
      cursor -> list_records(pds, did, collection, cursor, records)
    end
  end

  defp pds_from_did_document!(document) do
    document
    |> Map.get("service", [])
    |> Enum.find(fn service ->
      service["id"] == "#atproto_pds" or service["type"] == "AtprotoPersonalDataServer"
    end)
    |> case do
      %{"serviceEndpoint" => endpoint} when is_binary(endpoint) -> endpoint
      _ -> raise "could not find PDS serviceEndpoint in DID document"
    end
  end

  defp ensure_http_started! do
    {:ok, _apps} = Application.ensure_all_started(:inets)
    {:ok, _apps} = Application.ensure_all_started(:ssl)
  end

  defp get_json!(url) do
    ensure_http_started!()

    url
    |> String.to_charlist()
    |> :httpc.request()
    |> case do
      {:ok, {{_version, status, _reason}, _headers, body}} when status in 200..299 ->
        Jason.decode!(to_string(body))

      {:ok, {{_version, status, reason}, _headers, body}} ->
        raise "GET #{url} failed with #{status} #{reason}: #{to_string(body)}"

      {:error, reason} ->
        raise "GET #{url} failed: #{inspect(reason)}"
    end
  end
end
