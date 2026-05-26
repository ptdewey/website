defmodule Site.Pages.Post do
  use Phoenix.Component
  import Phoenix.HTML
  import Site.Components

  def routes(%{posts: posts}) do
    Enum.map(posts, fn post ->
      %Site.Route{
        path: post.path,
        title: post.title,
        page: page(%{post: post})
      }
    end)
  end

  defp page(assigns) do
    ~H"""
    <article class="post">
      <.post_header post={@post} />
      <%= raw @post.body %>
    </article>
    """
  end
end
