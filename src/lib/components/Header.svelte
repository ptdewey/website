<script lang="ts">
  import { base } from "$app/paths";

  import DarkModeToggle from "./DarkModeToggle.svelte";
</script>

<header>
  <nav class="flex prose prose-sm space-x-5">
    <a href="{base}/" class="no-underline">home</a>
    <a href="{base}/projects" class="no-underline">projects</a>
    <a href="{base}/writing" class="no-underline">writing</a>
    <DarkModeToggle />
  </nav>
</header>
