defmodule Mix.Tasks.SyncLeaflets do
  use Mix.Task

  @shortdoc "Fetch Leaflet site documents from an ATProto repo into content/leaflets.yml"

  @did "did:plc:hm5f3dnm6jdhrc55qp2npdja"
  @publication "at://did:plc:hm5f3dnm6jdhrc55qp2npdja/site.standard.publication/3mjdampndts24"
  @pds nil
  @base_url "https://coffee-thoughts.leaflet.pub"
  @output "content/leaflets.yml"

  @requirements ["app.start"]

  alias Site.StandardSite.Records

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
      |> Records.list()
      |> Records.filter_by_publication(publication)
      |> Enum.map(
        &Records.to_listing_entry(&1,
          base_url: base_url,
          source: "leaflet",
          default_title: "Untitled Leaflet"
        )
      )
      |> Enum.sort_by(& &1.date, {:desc, Date})

    File.mkdir_p!(Path.dirname(output))
    File.write!(output, to_yaml(leaflets))

    Mix.shell().info("Wrote #{length(leaflets)} leaflets to #{output}")
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
