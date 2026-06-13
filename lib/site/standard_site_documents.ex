defmodule Site.StandardSiteDocuments do
  @moduledoc false

  alias Site.StandardSite

  def all(path \\ StandardSite.documents_path()) do
    if File.exists?(path) do
      path
      |> YamlElixir.read_from_file!()
      |> case do
        nil -> %{}
        docs -> docs
      end
      |> Map.new(fn {post_path, attrs} -> {post_path, normalize(attrs)} end)
    else
      %{}
    end
  end

  def annotate_post(post, path \\ StandardSite.documents_path()) do
    metadata = Map.get(all(path), post.path, %{})

    %{
      post
      | standard_site_uri: Map.get(metadata, :uri) || post.standard_site_uri,
        standard_site_cid: Map.get(metadata, :cid) || post.standard_site_cid,
        bsky_post_ref: Map.get(metadata, :bsky_post_ref) || post.bsky_post_ref
    }
  end

  def write!(documents, path \\ StandardSite.documents_path()) when is_map(documents) do
    File.mkdir_p!(Path.dirname(path))

    documents
    |> Enum.sort_by(fn {post_path, _attrs} -> post_path end)
    |> Enum.map_join("", fn {post_path, attrs} -> document_yaml(post_path, attrs) end)
    |> then(fn yaml -> if yaml == "", do: "{}\n", else: yaml end)
    |> then(&File.write!(path, &1))
  end

  def merge_records(records, existing \\ %{}) do
    Enum.reduce(records, existing, fn record, acc ->
      case document_from_record(record) do
        nil -> acc
        {path, attrs} -> Map.put(acc, path, Map.merge(Map.get(acc, path, %{}), attrs))
      end
    end)
  end

  def document_from_record(%{"uri" => uri, "cid" => cid, "value" => value}) do
    with path when is_binary(path) <- value["path"] do
      attrs = %{
        uri: uri,
        cid: cid,
        rkey: StandardSite.rkey_from_uri(uri),
        title: value["title"],
        published_at: value["publishedAt"],
        updated_at: value["updatedAt"],
        bsky_post_ref: value["bskyPostRef"]
      }

      {path, reject_nil(attrs)}
    else
      _ -> nil
    end
  end

  defp normalize(attrs) when is_map(attrs) do
    %{
      uri: attrs["uri"],
      cid: attrs["cid"],
      rkey: attrs["rkey"],
      title: attrs["title"],
      published_at: attrs["published_at"],
      updated_at: attrs["updated_at"],
      bsky_post_ref: attrs["bsky_post_ref"]
    }
    |> reject_nil()
  end

  defp reject_nil(map) do
    Map.reject(map, fn {_key, value} -> is_nil(value) end)
  end

  defp document_yaml(post_path, attrs) do
    lines = [
      "#{yaml_string(post_path)}:",
      "  uri: #{yaml_string(attrs[:uri])}",
      "  cid: #{yaml_string(attrs[:cid])}",
      optional_line("  rkey", attrs[:rkey]),
      optional_line("  title", attrs[:title]),
      optional_line("  published_at", attrs[:published_at]),
      optional_line("  updated_at", attrs[:updated_at])
    ]

    lines =
      case attrs[:bsky_post_ref] do
        %{"uri" => uri, "cid" => cid} ->
          lines ++
            ["  bsky_post_ref:", "    uri: #{yaml_string(uri)}", "    cid: #{yaml_string(cid)}"]

        _ ->
          lines
      end

    lines
    |> Enum.reject(&is_nil/1)
    |> Enum.join("\n")
    |> Kernel.<>("\n")
  end

  defp optional_line(_key, nil), do: nil
  defp optional_line(key, value), do: "#{key}: #{yaml_string(value)}"

  defp yaml_string(nil), do: "\"\""

  defp yaml_string(value) do
    value
    |> to_string()
    |> String.replace("\\", "\\\\")
    |> String.replace("\"", "\\\"")
    |> then(&"\"#{&1}\"")
  end
end
