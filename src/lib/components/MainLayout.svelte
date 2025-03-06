<script>
  import Header from "./Header.svelte";
  export let data;
</script>

<div class="layout">
  <Header pages={data.pages} />
  <main>
    <slot />
  </main>
</div>

<style>
  .layout {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
  }
</style>
