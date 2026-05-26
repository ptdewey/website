defmodule Site.MarkdownConverter do
  def convert(_filepath, body, attrs, _opts) do
    Site.Markdown.to_html!(body, attrs)
  end
end
