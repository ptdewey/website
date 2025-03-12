<script lang="ts">
  import { base } from "$app/paths";
  import { onMount } from "svelte";

  onMount(() => {
    window.location.href = `${base}/writing/bluesky-comments-svelte`;
  });
</script>
