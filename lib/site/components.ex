defmodule Site.Components do
  use Phoenix.Component

  def header(assigns) do
    assigns = assign(assigns, :links, Site.Routes.nav_links())

    ~H"""
    <header>
      <nav>
        <ul>
          <li :for={link <- @links}><a href={link.path}>{link.name}</a></li>
        </ul>
      </nav>
    </header>
    """
  end

  def posts_list(assigns) do
    ~H"""
    <ul class="posts-list">
      <li :for={post <- @posts} class="posts-list__item">
        <time datetime={Date.to_iso8601(post.date)}><%= short_date(post.date) %></time>
        <a href={post.path}><%= post.title %></a>
      </li>
    </ul>
    """
  end

  def post_header(assigns) do
    ~H"""
    <header class="post-header">
      <h1><%= @post.title %></h1>
      <div class="post-meta">
        <time datetime={Date.to_iso8601(@post.date)}><%= long_date(@post.date) %></time>
        <span><%= @post.read_time %> min read</span>
        <ul class="tags" aria-label="tags">
          <li :for={tag <- @post.tags}><%= tag %></li>
        </ul>
      </div>
    </header>
    """
  end

  defp short_date(date), do: Calendar.strftime(date, "%Y-%m-%d")

  defp long_date(date), do: Calendar.strftime(date, "%B %-d, %Y")

  # TODO: implement components to use in md heex here
end
