defmodule Site.Generator do
  @output_dir "./output"
  @standard_site_publication "at://did:plc:hm5f3dnm6jdhrc55qp2npdja/site.standard.publication/3mgj4qfasw32n"

  def build() do
    File.rm_rf!(@output_dir)
    File.mkdir_p!(@output_dir)
    File.cp_r!("assets", Path.join([@output_dir, "assets"]))
    render_standard_site_well_known()

    Site.Blog.all_posts()
    |> Site.Routes.routes()
    |> Enum.each(&render_route/1)
  end

  defp render_route(route) do
    rendered =
      Site.Layouts.Root.render(%{
        title: route.title,
        content: route.page,
        standard_site_document: route.standard_site_document,
        standard_site_publication: route.standard_site_publication
      })

    route.path
    |> output_path()
    |> render_file(rendered)
  end

  defp output_path(path) do
    path =
      path
      |> String.trim_leading("/")
      |> String.trim_trailing("/")

    Path.join([@output_dir, path, "index.html"])
  end

  defp render_file(output, rendered) do
    File.mkdir_p!(Path.dirname(output))

    rendered
    |> Phoenix.HTML.Safe.to_iodata()
    |> IO.iodata_to_binary()
    |> then(&File.write!(output, &1))
  end

  defp render_standard_site_well_known do
    output = Path.join([@output_dir, ".well-known", "site.standard.publication"])

    File.mkdir_p!(Path.dirname(output))
    File.write!(output, @standard_site_publication)
  end
end
