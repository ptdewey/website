<script lang="ts">
  import { base } from "$app/paths";
  import type { Writing } from "$lib/types";
  import { formatDate } from "$lib/utils/dates";
  import { CommentSection } from "bluesky-comments-svelte";

  let data = $props();
  let writing: Writing = data.data;

  const author = "pdewey.com";
</script>

<link rel="stylesheet" href="{base}/darkearth-syntax.css" />

<svelte:head>
  <title>{writing.metadata.title}</title>
</svelte:head>

<article>
  {#if writing}
    <h2>{writing.metadata.title}</h2>
    <ul class="list-none list-outside pl-0 mb-0">
      <li class="pl-0">{formatDate(writing.metadata.date)}</li>
      <!-- {#if writing.metadata.categories && writing.metadata.categories.length > 1} -->
      <!--   <li class="pl-0"> -->
      <!--     Categories: {writing.metadata.categories.join(", ")} -->
      <!--   </li> -->
      <!-- {:else if writing.metadata.categories && writing.metadata.categories.length == 1} -->
      <!--   <li class="pl-0">Category: {writing.metadata.categories}</li> -->
      <!-- {/if} -->
      <li class="pl-0">{writing.metadata.read_time} minute read</li>
    </ul>
    <div class="break-words hyphens-auto">
      {@html writing.content}
    </div>

    {#if writing.metadata.tags}
      <hr />
      <ul class="list-none list-outside pl-0">
        {#if writing.metadata.tags.length > 1}
          <li class="pl-0">
            Tags: {writing.metadata.tags.join(", ")}
          </li>
        {:else}
          <li class="pl-0">Tag: {writing.metadata.tags}</li>
        {/if}
      </ul>
    {/if}

    <h2>Comments</h2>
    <div>
      {#if writing.metadata.bluesky_link}
        <CommentSection
          uri={writing.metadata.bluesky_link}
          opts={{ showCommentsTitle: false }}
        />
      {:else}
        <CommentSection {author} opts={{ showCommentsTitle: false }} />
      {/if}
    </div>
  {:else}
    <p>Post not found.</p>
  {/if}
</article>
