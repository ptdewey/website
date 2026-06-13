defmodule Site.Route do
  @enforce_keys [:path, :page]
  defstruct([
    :path,
    :page,
    :title,
    :standard_site_document,
    :standard_site_publication,
    nav?: false
  ])
end
