/**
 * MIT License
 *
 * Copyright (c) 2025 solanaceae.net
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 **/

function BlueskyComments(container, options = {}) {
  const css = `
.bsky-comment-section-container {
  font-family:
    system-ui,
    -apple-system,
    sans-serif;
}

.bsky-comments-list {
  margin-top: 0.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.bsky-error-text,
.bsky-loading-text {
  text-align: center;
}

.bsky-divider, .prose hr {
  margin: 0.5rem;
}

.bsky-comment-text-a,
.bsky-comment-text-p {
  margin-bottom: 0;
  color: var(--post-content-color);
}

.bsky-show-more-button {
  margin-top: 0.5rem;
  font-size: 0.875rem;
  color: var(--show-more-button-color, #3b82f6);
  text-decoration: underline;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
}

.bsky-stats-bar {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin: 0;
  margin-top: 0.5rem;
}

.bsky-stat-item {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  white-space: nowrap;
}

.bsky-stats-a {
  text-decoration: none;
  color: inherit;
}

.bsky-stat-icon {
  width: 1.25rem;
  height: 1.25rem;
}

.bsky-comments-title {
  margin-top: 1.5rem;
  margin-left: var(--left-margin-comments-title, 0);
  font-size: var(--font-size-title, 1.25rem);
  font-weight: bold;
}

.bsky-reply-text {
  margin-top: 0.5rem;
  font-size: var(--font-size-comment-body, 0.875rem);
}

.bsky-replies-container {
  border-left: 2px solid var(--comment-border-color, #525252);
  padding-left: 0.5rem;
}

.bsky-comment-container {
  margin: 0.5rem 0;
  font-size: 0.875rem;
}

.bsky-comment-content {
  margin-top: 0.35rem;
  display: flex;
  max-width: 36rem;
  flex-direction: column;
  align-items: var(--comment-content-alignment, flex-start);
}

.bsky-author-link {
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  align-items: center;
  color: var(--author-link-color);
  text-decoration: none;
}

.bsky-author-link:hover {
  text-decoration: underline;
}

.bsky-comment-avatar {
  height: var(--avatar-size, 1rem);
  width: var(--avatar-size, 1rem);
  flex-shrink: 0;
  border-radius: 9999px;
  background-color: var(--avatar-background-color, #d1d5db);
  margin: 0;
  margin-right: 0.3rem;
}

.bsky-author-name {
  overflow: hidden;
  line-clamp: 1;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  color: var(--author-link-color);
  margin: 0;
}

.bsky-comment-container a {
  text-decoration: none;
  color: inherit;
}

.bsky-comment-container a:hover {
  text-decoration: none;
}

.bsky-comment-content .bsky-handle {
  color: var(--handle-color, #6b7280);
}

.bsky-actions-icon {
  width: 1.25rem;
  height: 1.25rem;
}

.bsky-actions-container {
  margin-top: 0.35rem;
  display: flex;
  width: 100%;
  max-width: 150px;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  opacity: 0.6;
}

.bsky-actions-row {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.bsky-text-xs {
  font-size: 0.75rem;
  margin: 0;
}
`;

  function injectStyles() {
    if (typeof document === "undefined") return;

    const existingStyle = document.getElementById("bluesky-comments-styles");
    if (existingStyle) return;

    const style = document.createElement("style");
    style.id = "bluesky-comments-styles";
    style.textContent = css;
    document.head.appendChild(style);
  }

  const el =
    typeof container === "string"
      ? document.querySelector(container)
      : container;
  const opts = { uri: "", author: "", showCommentsTitle: true, ...options };

  const icon = (type, fill = "none", stroke = "currentColor") => {
    const paths = {
      reply:
        "M12 20.25c4.97 0 9-3.694 9-8.25s-4.03-8.25-9-8.25S3 7.444 3 12c0 2.104.859 4.023 2.273 5.48.432.447.74 1.04.586 1.641a4.483 4.483 0 0 1-.923 1.785A5.969 5.969 0 0 0 6 21c1.282 0 2.47-.402 3.445-1.087.81.22 1.668.337 2.555.337Z",
      repost:
        "M19.5 12c0-1.232-.046-2.453-.138-3.662a4.006 4.006 0 0 0-3.7-3.7 48.678 48.678 0 0 0-7.324 0 4.006 4.006 0 0 0-3.7 3.7c-.017.22-.032.441-.046.662M19.5 12l3-3m-3 3-3-3m-12 3c0 1.232.046 2.453.138 3.662a4.006 4.006 0 0 0 3.7 3.7 48.656 48.656 0 0 0 7.324 0 4.006 4.006 0 0 0 3.7-3.7c.017-.22.032-.441.046-.662M4.5 12l3 3m-3-3-3 3",
      like: "M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z",
    };
    return `<svg class="bsky-actions-icon" xmlns="http://www.w3.org/2000/svg" fill="${fill}" viewBox="0 0 24 24" stroke-width="1.5" stroke="${stroke}"><path stroke-linecap="round" stroke-linejoin="round" d="${paths[type]}" /></svg>`;
  };

  const stat = (type, count) =>
    `<div class="bsky-actions-row">${icon(type)}<p class="bsky-text-xs">${count || 0}</p></div>`;

  const actions = (post) => `
    <div class="bsky-actions-container">
      ${stat("reply", post.replyCount)}
      ${stat("repost", post.repostCount)}
      ${stat("like", post.likeCount)}
    </div>
  `;

  const sortReplies = (replies) =>
    replies
      .filter((r) => r.$type === "app.bsky.feed.defs#threadViewPost")
      .sort((a, b) => (b.post?.likeCount || 0) - (a.post?.likeCount || 0));

  const renderComment = (comment) => {
    if (
      !comment.post?.record ||
      comment.post.record.$type !== "app.bsky.feed.post"
    )
      return "";

    const { author, uri, record } = comment.post;
    const postRkey = uri.split("/").pop();
    const profileUrl = `https://bsky.app/profile/${author.did}`;
    const replies = comment.replies?.length
      ? `<div class="bsky-replies-container">${sortReplies(comment.replies).map(renderComment).join("")}</div>`
      : "";

    return `
      <div class="bsky-comment-container">
        <div class="bsky-comment-content">
          <a class="bsky-author-link" href="${profileUrl}" target="_blank" rel="noreferrer noopener">
            <img class="bsky-comment-avatar" src="${author.avatar || ""}" alt="" />
            <p class="bsky-author-name">
              ${author.displayName || author.handle}
              <span class="bsky-handle">@${author.handle}</span>
            </p>
          </a>
          <a class="bsky-comment-text-a" href="${profileUrl}/post/${postRkey}" target="_blank" rel="noreferrer noopener">
            <p class="bsky-comment-text-p">${record.text}</p>
          </a>
          ${actions(comment.post)}
        </div>
        ${replies}
      </div>
    `;
  };

  const render = (thread, postUrl, visibleCount = 5) => {
    const title = opts.showCommentsTitle
      ? '<h2 class="bsky-comments-title">Comments</h2>'
      : "";
    const sorted = sortReplies(thread.replies || []);
    const visible = sorted.slice(0, visibleCount);
    const showMoreBtn =
      visibleCount < sorted.length
        ? '<button class="bsky-show-more-button">Show more comments</button>'
        : "";

    const stats = [
      {
        type: "like",
        count: thread.post.likeCount,
        fill: "pink",
        stroke: "pink",
        label: "likes",
      },
      {
        type: "repost",
        count: thread.post.repostCount,
        fill: "none",
        stroke: "green",
        label: "reposts",
      },
      {
        type: "reply",
        count: thread.post.replyCount,
        fill: "#7FBADC",
        stroke: "#7FBADC",
        label: "replies",
      },
    ];

    el.innerHTML = `
      <div class="bsky-comment-section-container">
        <a class="bsky-stats-a" href="${postUrl}" target="_blank" rel="noreferrer noopener">
          <p class="bsky-stats-bar">
            ${stats
              .map(
                (s) => `
              <span class="bsky-stat-item">
                ${icon(s.type, s.fill, s.stroke)}
                <span class="bsky-stat-item-text">${s.count || 0} ${s.label}</span>
              </span>
            `,
              )
              .join("")}
          </p>
        </a>
        ${title}
        <p class="bsky-reply-text">
          Reply on Bluesky <a href="${postUrl}" target="_blank" rel="noreferrer noopener">here</a> to join the conversation.
        </p>
        <hr class="bsky-divider" />
        <div class="bsky-comments-list">
          ${visible.map(renderComment).join("")}
          ${showMoreBtn}
        </div>
      </div>
    `;

    const btn = el.querySelector(".bsky-show-more-button");
    if (btn)
      btn.addEventListener("click", () =>
        render(thread, postUrl, visibleCount + 5),
      );
  };

  const validateUri = async (uri) => {
    if (uri.startsWith("at://")) {
      const [, , did, _, rkey] = uri.split("/");
      return { uri, postUrl: `https://bsky.app/profile/${did}/post/${rkey}` };
    } else if (uri.includes("bsky.app/profile/")) {
      const match = uri.match(/profile\/([\w.:]+)\/post\/([\w]+)/);
      if (!match) throw new Error("Invalid URI format");

      const [, handleOrDid, postId] = match;
      let did = handleOrDid;

      if (!handleOrDid.startsWith("did:")) {
        const res = await fetch(
          `https://api.bsky.app/xrpc/com.atproto.identity.resolveHandle?handle=${handleOrDid}`,
        );
        const data = await res.json();
        did = data.did;
      }

      return {
        uri: `at://${did}/app.bsky.feed.post/${postId}`,
        postUrl: `https://bsky.app/profile/${handleOrDid}/post/${postId}`,
      };
    }

    throw new Error("Invalid URI format");
  };

  const getThread = async (uri) => {
    const res = await fetch(
      `https://api.bsky.app/xrpc/app.bsky.feed.getPostThread?${new URLSearchParams({ uri })}`,
      {
        method: "GET",
        headers: { Accept: "application/json" },
        cache: "no-store",
      },
    );

    if (!res.ok) throw new Error("Failed to fetch post thread");

    const data = await res.json();
    if (
      !data.thread ||
      data.thread.$type !== "app.bsky.feed.defs#threadViewPost"
    ) {
      throw new Error("Could not find thread");
    }

    return data.thread;
  };

  const showError = (msg) => {
    el.innerHTML = `
      <div class="bsky-comment-section-container">
        ${opts.showCommentsTitle ? '<h2 class="bsky-comments-title">Comments</h2>' : ""}
        <p class="bsky-error-text">${msg}</p>
      </div>
    `;
  };

  const showLoading = () => {
    el.innerHTML = `
      <div class="bsky-comment-section-container">
        ${opts.showCommentsTitle ? '<h2 class="bsky-comments-title">Comments</h2>' : ""}
        <p class="bsky-loading-text">Loading comments...</p>
      </div>
    `;
  };

  const init = async () => {
    if (!el) return;

    showLoading();

    try {
      let uri, postUrl;

      if (opts.uri) {
        ({ uri, postUrl } = await validateUri(opts.uri));
      } else if (opts.author) {
        const currentUrl = window.location.href;
        const apiUrl = `https://api.bsky.app/xrpc/app.bsky.feed.searchPosts?q=*&url=${encodeURIComponent(currentUrl)}&author=${opts.author}`;

        const res = await fetch(apiUrl);
        const data = await res.json();

        if (!data.posts?.length) {
          showError(`No matching post found for this URL.`);
          return;
        }

        const post = data.posts[0];
        ({ uri, postUrl } = await validateUri(post.uri));
      } else {
        return;
      }

      const thread = await getThread(uri);
      render(thread, postUrl);
    } catch (err) {
      console.error(err);
      showError(err.message || "Error loading comments");
    }
  };

  injectStyles();
  init();
}

if (typeof module !== "undefined" && module.exports) {
  module.exports = BlueskyComments;
}
