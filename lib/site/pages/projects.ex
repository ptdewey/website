defmodule Site.Pages.Projects do
  use Phoenix.Component
  import Site.Components, warn: false

  @external_resource "content/pages/projects.md"
  @markdown File.read!("content/pages/projects.md")

  def route(_ctx) do
    %Site.Route{
      path: "/projects",
      title: "projects",
      nav?: true,
      page: page(%{})
    }
  end

  defp page(assigns) do
    Site.Markdown.to_heex!(@markdown, assigns)
  end
end
