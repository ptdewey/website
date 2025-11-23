<script lang="ts">
  import type { Writing } from "$lib/types";
  import { base } from "$app/paths";
  import { formatDate } from "$lib/utils/dates";

  let data = $props();
  let writings: Writing[] = data.data.data;
</script>

<svelte:head>
  <title>Writing</title>
</svelte:head>

<div>
  <h2>Writing</h2>
  <ul class="space-y-5 p-0">
    {#each writings as writing}
      <li class="flex w-full items-start">
        <span class="mr-5">{formatDate(writing.metadata.date)}</span>
        <div class="flex-3/4 overflow-hidden">
          <a href="{base}/writing/{writing.metadata.slug}"
            >{writing.metadata.title}</a
          >
          <p class="italic !m-0 p-0">{writing.metadata.type}</p>
        </div>
      </li>
    {/each}
  </ul>
</div>
