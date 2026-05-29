defmodule Site.Layouts.Root do
  use Phoenix.Component
  import Site.Components

  attr(:title, :string, default: "patrick dewey")
  attr(:content, :any, required: true)
  attr(:standard_site_document, :string, default: nil)
  attr(:standard_site_publication, :string, default: nil)

  # TODO: use page title in H1 if provided
  def render(assigns) do
    ~H"""
    <!DOCTYPE html>
    <html>
      <head>
        <meta charset="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <title>{@title}</title>
        <meta name="color-scheme" content="light" />
        <meta name="theme-color" content="#f0ebe1" />
        <style>
          :root {
            background: #f0ebe1;
            color: #252f1e;
          }

          html {
            background: #f0ebe1;
          }
        </style>
        <link rel="icon" href="/assets/favicon.ico" sizes="any" />
        <link rel="icon" type="image/png" href="/assets/favicon.png" />
        <link :if={@standard_site_document} rel="site.standard.document" href={@standard_site_document} />
        <link :if={@standard_site_publication} rel="site.standard.publication" href={@standard_site_publication} />
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
