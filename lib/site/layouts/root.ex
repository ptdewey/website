defmodule Site.Layouts.Root do
  use Phoenix.Component
  import Site.Components

  attr(:title, :string, default: "patrick dewey")
  attr(:content, :any, required: true)

  # TODO: use page title in H1 if provided
  def render(assigns) do
    ~H"""
    <!DOCTYPE html>
    <html>
      <head>
        <meta charset="utf-8" />
        <title>{@title}</title>
        <link rel="icon" href="/assets/favicon.ico" sizes="any" />
        <link rel="icon" type="image/png" href="/assets/favicon.png" />
        <link rel="stylesheet" href="/assets/css/app.css" />
      </head>
      <body>
        <.header />
        <main>{@content}</main>
      </body>
    </html>
    """
  end
end
