defmodule Site.Pages.Home do
  use Phoenix.Component

  @links [
    %{
      name: "atproto",
      text: "at://did:plc:hm5f3dnm6jdhrc55qp2npdja (@pdewey.com)",
      url: "https://pds.ls/at://did:plc:hm5f3dnm6jdhrc55qp2npdja"
    },
    %{
      name: "tangled",
      text: "@pdewey.com",
      url: "https://tangled.org/pdewey.com"
    },
    %{
      name: "github",
      text: "ptdewey",
      url: "https://github.com/ptdewey"
    }
  ]

  def route(%{posts: posts}) do
    %Site.Route{
      path: "/",
      title: "home",
      nav?: true,
      page: page(%{posts: posts})
    }
  end

  defp page(assigns) do
    assigns = Map.put(assigns, :links, @links)

    ~H"""
    <article>
      <h1>patrick dewey</h1>
      <section id="about">
        <p>
          Hi I'm Patrick, a software engineer interested in dev tools, 
          programming languages, coffee, music, and the open social web.
        </p>
        <p>
          I'm currently working on <a href="https://alpha.arabica.social">arabica.social</a>,
          a social coffee journaling site, built on <a href="https://atproto.com">atproto</a>.
        </p>
      </section>
      <section id="projects">
        <p>A few logs and lists live here too:</p>
        <ul>
          <li><a href="/projects">projects</a> - things I've made and contributed to</li>
          <li><a href="/uses">uses</a> - software, hardware, and other things I use</li>
          <li><a href="/coffee">coffee</a> - brews logged from arabica.social</li>
          <li><a href="/music">music</a> - listening history from teal.fm</li>
        </ul>
      </section>
      <section id="links">
        <p>You can find me around the internet at:</p>
        <ul>
          <li :for={link <- @links}>
            {link.name}: <a href={link.url}>{link.text}</a>
          </li>
        </ul>
      </section>
    </article>
    """
  end
end
