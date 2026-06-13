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
      The full journal, gear, and ratings are on <a href="https://alpha.arabica.social/profile/pdewey.com">alpha.arabica.social/profile/pdewey.com</a>.
    </p>
    <div id="brew-list" class="brew-list">
      <p class="log-msg">Fetching brews…</p>
    </div>

    <template id="brew-tmpl">
      <article class="feed-card feed-card-brew brew-card">
        <header class="brew-card__meta">
          <time class="brew-date"></time>
        </header>

        <div class="brew-card__body">
          <div class="brew-card__topline">
            <div class="brew-card__title-wrap">
              <a class="brew-bean" target="_blank" rel="noopener"></a>
              <div class="brew-roaster"></div>
            </div>
            <span class="brew-rating" aria-label="rating"></span>
          </div>

          <div class="brew-sub"></div>

          <div class="brew-brewer-row">
            <span class="brew-label">Brewer:</span>
            <span class="brew-brewer"></span>
          </div>

          <dl class="brew-params">
            <div class="brew-param brew-param-grinder">
              <dt>Grinder:</dt>
              <dd></dd>
            </div>
            <div class="brew-param brew-param-water">
              <dt>Water:</dt>
              <dd></dd>
            </div>
            <div class="brew-param brew-param-temp">
              <dt>Temp:</dt>
              <dd></dd>
            </div>
            <div class="brew-param brew-param-time">
              <dt>Time:</dt>
              <dd></dd>
            </div>
          </dl>

          <div class="brew-pours" aria-label="pours"></div>
          <p class="brew-notes"></p>
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
