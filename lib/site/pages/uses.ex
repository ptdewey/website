defmodule Site.Pages.Uses do
  use Phoenix.Component
  import Site.Components, warn: false

  @external_resource "content/pages/uses.md"
  @markdown File.read!("content/pages/uses.md")

  def route(_ctx) do
    %Site.Route{
      path: "/uses",
      title: "uses",
      nav?: false,
      page: page(%{})
    }
  end

  defp page(assigns) do
    Site.Markdown.to_heex!(@markdown, assigns)
  end
end
