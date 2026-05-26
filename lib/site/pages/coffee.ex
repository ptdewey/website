defmodule Site.Pages.Coffee do
  use Phoenix.Component
  import Site.Components, warn: false

  def route(_ctx) do
    %Site.Route{
      path: "/coffee",
      title: "coffee",
      page: page(%{})
    }
  end

  defp page(assigns) do
    ~H"""
    <h1>coffee log</h1>
    <p>
      Brews I've logged on <a href="https://alpha.arabica.social">arabica.social</a>, a coffee journal on the AT Protocol.
    </p>
    <p>
      The full journal, gear, and ratings are on <a href="https://alpha.arabica.social/profile/pdewey.com">@pdewey.com</a>.
    </p>
    <div id="brew-list" class="log-body">
      <p class="log-msg">Fetching brews…</p>
    </div>

    <template id="brew-tmpl">
      <article class="log-entry brew-entry">
        <time class="log-date"></time>
        <div class="log-value">
          <div class="log-head">
            <a class="brew-bean" target="_blank" rel="noopener"></a>
            <span class="brew-rating" aria-label="rating"></span>
          </div>
          <div class="brew-roaster"></div>
          <div class="brew-sub"></div>
          <div class="brew-meta"></div>
          <div class="brew-equipment"></div>
          <p class="brew-notes"></p>
          <div class="log-rkey brew-rkey"></div>
        </div>
      </article>
    </template>

    <noscript>
      <p>Enable JavaScript to view brews, or visit <a href="https://arabica.social">arabica.social</a>.</p>
    </noscript>

    <script type="module" src="/assets/js/coffee.js">
    </script>
    """
  end
end
