<script>
  import { CommentSection } from "bluesky-comments-svelte";
  import { base } from "$app/paths";

  /** @import {BlogPost} from "$lib/types" */

  /** @type {BlogPost} */
  export let data;

  const author = "pdewey.com";
</script>

<svelte:head>
  <title>{data.post.metadata.title}</title>
</svelte:head>

<link rel="stylesheet" href="{base}/darkearth-syntax.css" />
<article class="post-container">
  {#if data && data.post}
    <h1>{data.post.metadata.title}</h1>
    <ul class="post-metadata">
      <li>Date: {data.post.metadata.date}</li>
      {#if data.post.metadata.categories && data.post.metadata.categories.length > 1}
        <li>Categories: {data.post.metadata.categories.join(", ")}</li>
      {:else if data.post.metadata.categories && data.post.metadata.categories.length == 1}
        <li>Category: {data.post.metadata.categories}</li>
      {/if}
      <li>{data.post.metadata.read_time} minute read</li>
    </ul>
    <div class="post-content">
      {@html data.post.content}
    </div>

    {#if data.post.metadata.tags}
      <hr class="divider" />
      <ul class="post-metadata">
        {#if data.post.metadata.tags.length > 1}
          <li>Published with tags: {data.post.metadata.tags.join(", ")}</li>
        {:else}
          <li>Published with tag: {data.post.metadata.tags}</li>
        {/if}
      </ul>
    {/if}
  {:else}
    <p>Post not found.</p>
  {/if}

  <h2>Comments</h2>
  <div class="comment-section">
    <CommentSection {author} opts={{ showCommentsTitle: false }} />
  </div>
</article>

<style>
  :root {
    --comment-border-color: var(--tan);
    --avatar-size: 2rem;
    --font-size-title: 1.5rem;
    --handle-color: var(--green);
    --comment-content-alignment: flex-start;
    --font-size-comment-body: 1rem;
    --show-more-button-color: var(--tan);
  }

  .post-container {
    max-width: 800px;
  }

  .post-metadata {
    background-color: var(--background-color-alt);
    padding: 5px 15px;
    border-radius: 8px;
  }

  .post-metadata li {
    list-style: none;
    padding: 0px;
    margin: 0px;
  }

  .comment-section {
    background-color: var(--dark-brown-alt);
    padding: 10px;
    border-radius: 8px;
    max-width: 100%;
    box-sizing: border-box;
    overflow-x: hidden;
    overflow-wrap: break-word;
    white-space: normal;
  }

  .divider {
    color: var(--green);
  }
</style>
