defmodule Site.StandardSite.Records do
  @moduledoc false

  alias Site.StandardSite

  def list(opts) do
    opts
    |> Keyword.put(:collection, StandardSite.collection())
    |> Site.AtProto.list_records()
  end

  def filter_by_publication(records, nil), do: records

  def filter_by_publication(records, publication) when is_binary(publication) do
    Enum.filter(records, fn %{"value" => value} ->
      StandardSite.publication_uri(value) == publication
    end)
  end

  def to_listing_entry(%{"uri" => uri, "value" => value}, opts \\ []) do
    base_url = Keyword.fetch!(opts, :base_url)
    source = Keyword.get(opts, :source, "standard.site")
    default_title = Keyword.get(opts, :default_title, "Untitled Document")
    rkey = StandardSite.rkey_from_uri(uri)
    created_at = value["createdAt"] || value["publishedAt"] || value["updatedAt"]

    %{
      title: value["title"] || value["name"] || default_title,
      path:
        value["url"] || value["canonicalUrl"] || document_url(base_url, value["path"] || rkey),
      date: date!(created_at),
      source: source,
      at_uri: uri
    }
  end

  def document_url(base_url, path) do
    String.trim_trailing(base_url, "/") <> "/" <> String.trim_leading(path, "/")
  end

  def date!(nil),
    do: raise("Standard.site record did not include createdAt/publishedAt/updatedAt")

  def date!(datetime) do
    case DateTime.from_iso8601(datetime) do
      {:ok, datetime, _offset} -> DateTime.to_date(datetime)
      {:error, _reason} -> Date.from_iso8601!(datetime)
    end
  end
end
