defmodule Site.Pages.Music do
  use Phoenix.Component
  import Site.Components, warn: false

  def route(_ctx) do
    %Site.Route{
      path: "/music",
      title: "music",
      page: page(%{})
    }
  end

  defp page(assigns) do
    ~H"""
    <h1>listening log</h1>
    <p>
      Plays scrobbled by <a href="https://github.com/teal-fm/piper">piper</a> and written to my PDS for <a href="https://teal.fm">teal.fm</a>.
    </p>

    <div id="now-playing" class="now-playing hidden" aria-live="polite" aria-atomic="false">
      <span class="np-label">now playing</span>
      <div class="log-value">
        <div>
          <span id="np-artist" class="np-artist"></span>
          <span class="np-sep" aria-hidden="true">-</span>
          <span id="np-track" class="np-track"></span>
        </div>
        <div id="np-release" class="np-release"></div>
      </div>
    </div>

    <div id="play-list" class="log-body">
      <p class="log-msg">Fetching plays…</p>
    </div>

    <template id="play-tmpl">
      <article class="log-entry play-entry">
        <time class="log-date"></time>
        <div class="log-value">
          <a class="play-track" target="_blank" rel="noopener"></a>
          <div class="play-artist"></div>
          <div class="play-release"></div>
          <div class="log-rkey play-rkey"></div>
        </div>
      </article>
    </template>

    <noscript>
      <p>Enable JavaScript to view plays, or visit <a href="https://teal.fm">teal.fm</a>.</p>
    </noscript>

    <script type="module" src="/assets/js/music.js">
    </script>
    """
  end
end
