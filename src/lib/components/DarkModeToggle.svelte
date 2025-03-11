<script lang="ts">
  let isDark = false;

  if (typeof document !== "undefined") {
    isDark = document.documentElement.getAttribute("data-theme") === "dark";
  }

  function toggleTheme() {
    isDark = !isDark;
    document.documentElement.setAttribute(
      "data-theme",
      isDark ? "dark" : "light",
    );
    document.documentElement.classList.toggle("dark", isDark);
    localStorage.setItem("theme", isDark ? "dark" : "light");
  }
</script>

<button
  on:click={toggleTheme}
  class="prose space-x-5 ml-auto mr-0"
  aria-label="Toggle Dark Mode"
>
  {#if isDark}
    light
  {:else}
    dark
  {/if}
</button>
