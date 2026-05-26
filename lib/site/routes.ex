defmodule Site.Routes do
  @pages [
    Site.Pages.Home,
    Site.Pages.Projects,
    Site.Pages.Coffee,
    Site.Pages.Music,
    Site.Pages.Uses,
    Site.Pages.Blog,
    Site.Pages.Now
  ]

  def routes(posts) do
    ctx = %{posts: posts}

    Enum.map(@pages, & &1.route(ctx)) ++ Site.Pages.Post.routes(ctx)
  end

  def nav_links(posts \\ []) do
    posts
    |> routes()
    |> Enum.filter(& &1.nav?)
    |> Enum.map(fn route ->
      %{name: route.title, path: route.path}
    end)
  end
end
