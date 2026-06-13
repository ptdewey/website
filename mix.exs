defmodule Site.MixProject do
  use Mix.Project

  def project do
    [
      app: :site,
      version: "0.1.0",
      elixir: "~> 1.18",
      start_permanent: Mix.env() == :prod,
      deps: deps()
    ]
  end

  # Run "mix help compile.app" to learn about applications.
  def application do
    [
      extra_applications: [:logger]
    ]
  end

  # Run "mix help deps" to learn about dependencies.
  defp deps do
    [
      {:nimble_publisher, "~> 1.0", runtime: false},
      {:phoenix_live_view, "~> 1.1"},
      {:mdex, "~> 0.11"},
      {:yaml_elixir, "~> 2.12"},
      {:jason, "~> 1.4"},
      {:atex, git: "https://tangled.org/comet.sh/atex", branch: "main"}
    ]
  end
end
