defmodule Site.Markdown do
  require MDEx

  @options [
    extension: [
      strikethrough: true,
      table: true,
      footnotes: true,
      autolink: true
    ],
    syntax_highlight: [
      formatter: {:html_inline, theme: "gruvbox_dark"}
    ]
  ]

  def to_html!(body, assigns \\ %{}) do
    MDEx.to_html!(body, Keyword.put(@options, :assigns, Map.new(assigns)))
  end

  def to_heex!(body, assigns \\ %{}) do
    MDEx.to_heex!(body, Keyword.put(@options, :assigns, Map.new(assigns)))
  end
end
