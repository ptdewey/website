defmodule Site.Pages.Now do
  use Phoenix.Component
  import Site.Components, warn: false

  @external_resource "content/pages/now.md"
  @markdown File.read!("content/pages/now.md")

  def route(_ctx) do
    %Site.Route{
      path: "/now",
      title: "now",
      nav?: true,
      page: page(%{})
    }
  end

  defp page(assigns) do
    Site.Markdown.to_heex!(@markdown, assigns)
  end
end
