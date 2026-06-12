defmodule Site.Pages.Blog do
  use Phoenix.Component
  import Site.Components

  @standard_site_publication "at://did:plc:hm5f3dnm6jdhrc55qp2npdja/site.standard.publication/3mgj4qfasw32n"

  def route(%{posts: posts}) do
    posts =
      posts
      |> Kernel.++(Site.Leaflets.all())
      |> Enum.sort_by(& &1.date, {:desc, Date})

    %Site.Route{
      path: "/blog",
      title: "blog",
      nav?: true,
      standard_site_publication: @standard_site_publication,
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
