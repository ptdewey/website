defmodule Site.StandardSite do
  @moduledoc false

  @did "did:plc:hm5f3dnm6jdhrc55qp2npdja"
  @publication "at://#{@did}/site.standard.publication/3mgj4qfasw32n"
  @documents_path "content/standard_site_documents.yml"
  @collection "site.standard.document"

  def did, do: @did
  def publication, do: @publication
  def documents_path, do: @documents_path
  def collection, do: @collection

  def rkey_from_uri(nil), do: nil

  def rkey_from_uri(uri) when is_binary(uri) do
    uri |> String.split("/") |> List.last()
  end

  def slug_rkey(%{id: id}) when is_binary(id) do
    if valid_rkey?(id), do: id, else: Atex.TID.now() |> to_string()
  end

  def document_path(%{id: id}) when is_binary(id), do: "/" <> id

  def valid_rkey?(rkey) when is_binary(rkey) do
    byte_size(rkey) <= 512 and Regex.match?(~r/^[A-Za-z0-9._:~-]+$/, rkey) and
      rkey not in [".", ".."]
  end

  def record_for_post(post, opts \\ []) do
    publication = Keyword.get(opts, :publication, @publication)
    published_at = datetime_for_date(post.date)
    updated_at = Keyword.get(opts, :updated_at) || Map.get(post, :updated_at)

    %{
      "$type" => @collection,
      "site" => publication,
      "path" => document_path(post),
      "title" => post.title,
      "publishedAt" => published_at,
      "content" => markdown_content(post),
      "textContent" => text_content(post.body)
    }
    |> put_if_present("updatedAt", normalize_datetime(updated_at))
    |> put_if_present("description", Map.get(post, :standard_site_description))
    |> put_if_present("tags", normalize_tags(Map.get(post, :tags)))
    |> put_if_present("bskyPostRef", normalize_strong_ref(Map.get(post, :bsky_post_ref)))
  end

  def publication_uri(%{"site" => site}) when is_binary(site), do: site

  def publication_uri(%{"publication" => publication}) when is_binary(publication),
    do: publication

  def publication_uri(%{"publication" => %{"uri" => publication}}), do: publication
  def publication_uri(_value), do: nil

  defp put_if_present(map, _key, nil), do: map
  defp put_if_present(map, _key, []), do: map
  defp put_if_present(map, key, value), do: Map.put(map, key, value)

  defp normalize_tags(nil), do: []
  defp normalize_tags(tags) when is_list(tags), do: Enum.reject(tags, &is_nil/1)
  defp normalize_tags(_tags), do: []

  defp normalize_strong_ref(nil), do: nil

  defp normalize_strong_ref(%{"uri" => uri, "cid" => cid}) when is_binary(uri) and is_binary(cid),
    do: %{"uri" => uri, "cid" => cid}

  defp normalize_strong_ref(%{uri: uri, cid: cid}) when is_binary(uri) and is_binary(cid),
    do: %{"uri" => uri, "cid" => cid}

  defp normalize_strong_ref(_ref), do: nil

  defp markdown_content(%{markdown: markdown}) when is_binary(markdown) and markdown != "" do
    %{
      "$type" => "com.pdewey.document.markdown",
      "markdown" => markdown
    }
  end

  defp markdown_content(_post), do: nil

  defp datetime_for_date(%Date{} = date) do
    date
    |> DateTime.new!(~T[00:00:00], "Etc/UTC")
    |> DateTime.to_iso8601()
  end

  defp normalize_datetime(nil), do: nil
  defp normalize_datetime(%DateTime{} = datetime), do: DateTime.to_iso8601(datetime)
  defp normalize_datetime(%Date{} = date), do: datetime_for_date(date)
  defp normalize_datetime(datetime) when is_binary(datetime), do: datetime

  defp text_content(html) do
    html
    |> String.replace(~r/<(script|style)[^>]*>.*?<\/\1>/si, " ")
    |> String.replace(~r/<[^>]+>/, " ")
    |> html_unescape()
    |> String.replace(~r/\s+/, " ")
    |> String.trim()
  end

  defp html_unescape(text) do
    text
    |> String.replace("&amp;", "&")
    |> String.replace("&lt;", "<")
    |> String.replace("&gt;", ">")
    |> String.replace("&quot;", "\"")
    |> String.replace("&#39;", "'")
  end
end
