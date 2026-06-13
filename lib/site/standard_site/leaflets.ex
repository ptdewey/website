defmodule Site.StandardSite.Leaflets do
  @moduledoc false

  @path "content/leaflets.yml"

  def all(path \\ @path) do
    if File.exists?(path) do
      path
      |> YamlElixir.read_from_file!()
      |> Enum.map(&normalize/1)
    else
      []
    end
  end

  defp normalize(attrs) do
    %{
      title: Map.fetch!(attrs, "title"),
      path: Map.fetch!(attrs, "path"),
      date: attrs |> Map.fetch!("date") |> Date.from_iso8601!(),
      source: Map.get(attrs, "source", "leaflet"),
      at_uri: Map.get(attrs, "at_uri")
    }
  end
end
