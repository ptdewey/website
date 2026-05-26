defmodule Site.Generator do
  @output_dir "./output"

  def build() do
    File.rm_rf!(@output_dir)
    File.mkdir_p!(@output_dir)
    File.cp_r!("assets", Path.join([@output_dir, "assets"]))

    Site.Blog.all_posts()
    |> Site.Routes.routes()
    |> Enum.each(&render_route/1)
  end

  defp render_route(route) do
    rendered =
      Site.Layouts.Root.render(%{
        title: route.title,
        content: route.page
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
end
