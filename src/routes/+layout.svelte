<script lang="ts">
  import Header from "$lib/components/Header.svelte";
  import "../app.css";
  let { children } = $props();
</script>

<!-- <div class="font-display flex justify-center min-h-screen p-6 bg-[#f4f6f2]"> -->
<div
  class="font-display flex justify-center min-h-screen p-6 bg-background text-primary"
>
  <div
    class="prose prose-a:no-underline w-full max-w-prose text-[16px] leading-relaxed mt-6"
  >
    <Header />
    <main class="link-color prose prose-sm">
      {@render children()}
    </main>
  </div>
</div>
