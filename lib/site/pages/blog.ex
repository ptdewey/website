defmodule Site.Pages.Blog do
  use Phoenix.Component
  import Site.Components

  def route(%{posts: posts}) do
    %Site.Route{
      path: "/blog",
      title: "blog",
      nav?: true,
      page: page(%{posts: posts})
    }
  end

  defp page(assigns) do
    ~H"""
    <article>
      <h1>blog</h1>
      <.posts_list posts={@posts} />
    </article>
    """
  end
end
