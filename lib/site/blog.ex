defmodule Site.Blog do
  alias Site.Post

  # Parse and convert markdown to HTML at compile time
  use NimblePublisher,
    build: Post,
    from: "./content/posts/**/*.md",
    as: :posts,
    parser: Site.FrontmatterParser,
    html_converter: Site.MarkdownConverter

  @posts Enum.sort_by(@posts, & &1.date, {:desc, Date})

  def all_posts, do: @posts
end
