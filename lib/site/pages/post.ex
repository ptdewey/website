defmodule Site.Pages.Post do
  use Phoenix.Component
  import Phoenix.HTML
  import Site.Components

  @standard_site_publication "at://did:plc:hm5f3dnm6jdhrc55qp2npdja/site.standard.publication/3mgj4qfasw32n"

  def routes(%{posts: posts}) do
    Enum.map(posts, fn post ->
      post = Site.StandardSite.Documents.annotate_post(post)

      %Site.Route{
        path: post.path,
        title: post.title,
        standard_site_document: post.standard_site_uri,
        standard_site_publication: @standard_site_publication,
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
