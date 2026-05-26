defmodule Site.Route do
  @enforce_keys [:path, :page]
  defstruct([:path, :page, :title, nav?: false])
end
