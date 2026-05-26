defmodule Site.FrontmatterParser do
  @known_keys %{
    "title" => :title,
    "author" => :author,
    "tags" => :tags,
    "description" => :description
  }

  def parse(_path, contents) do
    ["---\n" <> yaml, body] = :binary.split(contents, "\n---\n")

    attrs =
      yaml
      |> YamlElixir.read_from_string!()
      |> atomize_known_keys()

    {attrs, body}
  end

  defp atomize_known_keys(attrs) do
    Map.new(attrs, fn {key, value} ->
      {Map.fetch!(@known_keys, key), value}
    end)
  end
end
