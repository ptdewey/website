<script lang="ts">
  import Footer from "$lib/components/Footer.svelte";
  import Header from "$lib/components/Header.svelte";
  import "../app.css";
  let { children } = $props();
</script>

<!-- <div class="font-display flex justify-center min-h-screen p-6 bg-[#f4f6f2]"> -->
<div
  class="font-display min-h-screen flex flex-col justify-between bg-background text-primary p-6"
>
  <div
    class="w-full max-w-prose mx-auto flex-grow prose prose-a:no-underline text-[16px] leading-relaxed"
  >
    <Header />
    <main class="link-color prose prose-sm">
      {@render children()}
    </main>

    <Footer />
  </div>
</div>
