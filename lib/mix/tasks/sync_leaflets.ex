defmodule Mix.Tasks.SyncLeaflets do
  use Mix.Task

  @shortdoc "Fetch Leaflet site documents from an ATProto repo into content/leaflets.yml"

  @did "did:plc:hm5f3dnm6jdhrc55qp2npdja"
  @publication "at://did:plc:hm5f3dnm6jdhrc55qp2npdja/site.standard.publication/3mjdampndts24"
  @pds nil
  @base_url "https://coffee-thoughts.leaflet.pub"
  @output "content/leaflets.yml"

  @requirements ["app.start"]

  @impl Mix.Task
  def run(args) do
    {opts, _args, _invalid} =
      OptionParser.parse(args,
        strict: [
          did: :string,
          publication: :string,
          pds: :string,
          base_url: :string,
          output: :string
        ]
      )

    did =
      opts[:did] || @did ||
        raise "pass --did, or set @did in #{__ENV__.file}"

    publication =
      opts[:publication] || @publication

    pds = opts[:pds] || @pds
    base_url = opts[:base_url] || @base_url
    output = opts[:output] || @output

    leaflets =
      [did: did, pds: pds]
      |> Site.AtProto.list_site_documents()
      |> maybe_filter_publication(publication)
      |> Enum.map(&to_leaflet(&1, base_url))
      |> Enum.sort_by(& &1.date, {:desc, Date})

    File.mkdir_p!(Path.dirname(output))
    File.write!(output, to_yaml(leaflets))

    Mix.shell().info("Wrote #{length(leaflets)} leaflets to #{output}")
  end

  defp maybe_filter_publication(records, nil), do: records

  defp maybe_filter_publication(records, publication) do
    Enum.filter(records, fn %{"value" => value} -> publication_uri(value) == publication end)
  end

  defp publication_uri(%{"site" => site}) when is_binary(site), do: site

  defp publication_uri(%{"publication" => publication}) when is_binary(publication),
    do: publication

  defp publication_uri(%{"publication" => %{"uri" => publication}}), do: publication
  defp publication_uri(_value), do: nil

  defp to_leaflet(%{"uri" => uri, "value" => value}, base_url) do
    rkey = uri |> String.split("/") |> List.last()
    created_at = value["createdAt"] || value["publishedAt"] || value["updatedAt"]

    %{
      title: value["title"] || value["name"] || "Untitled Leaflet",
      path:
        value["url"] || value["canonicalUrl"] || document_url(base_url, value["path"] || rkey),
      date: date!(created_at),
      source: "leaflet",
      at_uri: uri
    }
  end

  defp document_url(base_url, path) do
    String.trim_trailing(base_url, "/") <> "/" <> String.trim_leading(path, "/")
  end

  defp date!(nil), do: raise("leaflet record did not include createdAt/publishedAt/updatedAt")

  defp date!(datetime) do
    case DateTime.from_iso8601(datetime) do
      {:ok, datetime, _offset} -> DateTime.to_date(datetime)
      {:error, _reason} -> Date.from_iso8601!(datetime)
    end
  end

  defp to_yaml([]), do: "[]\n"

  defp to_yaml(leaflets) do
    Enum.map_join(leaflets, "", fn leaflet ->
      """
      - title: #{yaml_string(leaflet.title)}
        path: #{yaml_string(leaflet.path)}
        date: #{yaml_string(Date.to_iso8601(leaflet.date))}
        source: #{yaml_string(leaflet.source)}
        at_uri: #{yaml_string(leaflet.at_uri)}
      """
    end)
  end

  defp yaml_string(value) do
    value
    |> to_string()
    |> String.replace("\\", "\\\\")
    |> String.replace("\"", "\\\"")
    |> then(&"\"#{&1}\"")
  end
end
