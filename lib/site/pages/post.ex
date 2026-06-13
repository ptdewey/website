defmodule Site.Pages.Post do
  use Phoenix.Component
  import Phoenix.HTML
  import Site.Components

  @standard_site_publication "at://did:plc:hm5f3dnm6jdhrc55qp2npdja/site.standard.publication/3mgj4qfasw32n"

  def routes(%{posts: posts}) do
    Enum.map(posts, fn post ->
      post = Site.StandardSite.Documents.annotate_post(post)

      %Site.Route{
        path: post.path,
        title: post.title,
        standard_site_document: post.standard_site_uri,
        standard_site_publication: @standard_site_publication,
        page: page(%{post: post, bsky_conversation_url: bsky_conversation_url(post)})
      }
    end)
  end

  defp page(assigns) do
    ~H"""
    <article class="post">
      <.post_header post={@post} />
      <%= raw @post.body %>
    </article>

    <section :if={@bsky_conversation_url} class="bsky-comments" aria-labelledby="bsky-comments-title">
      <script type="module" src="/assets/js/bsky-conversation.js">
      </script>
      <h2 id="bsky-comments-title">comments</h2>
      <bsky-conversation
        uri={@bsky_conversation_url}
        max-depth="3"
        engage-text="add your thoughts on bluesky"
        header-template="{replies?{replies|reply|replies}}{quotes?, {quotes|quote|quotes}}{repostedBy?, reposted by {repostedBy}}"
      >
      </bsky-conversation>
    </section>
    """
  end

  defp bsky_conversation_url(%{bsky_post_ref: %{"uri" => uri}}), do: bsky_post_url(uri)
  defp bsky_conversation_url(%{bsky_post_ref: %{uri: uri}}), do: bsky_post_url(uri)

  defp bsky_conversation_url(%{bsky_post_ref: uri}) when is_binary(uri),
    do: bsky_post_url(uri)

  defp bsky_conversation_url(_post), do: nil

  defp bsky_post_url("at://" <> rest) do
    case String.split(rest, "/", parts: 3) do
      [repo, "app.bsky.feed.post", rkey] -> "https://bsky.app/profile/#{repo}/post/#{rkey}"
      _ -> nil
    end
  end

  defp bsky_post_url("https://bsky.app/profile/" <> _ = url), do: url
  defp bsky_post_url(_uri), do: nil
end
