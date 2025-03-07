<script>
  import { base } from "$app/paths";
  /** @import {Page} from "$lib/types" */

  /** @type {Page[]} */
  export let pages = [];
  let isMenuOpen = false;

  function toggleMenu() {
    isMenuOpen = !isMenuOpen;
  }

  function closeMenu() {
    isMenuOpen = false;
  }
</script>

<header>
  <nav>
    <a href="{base}/" class="logo">Home</a>
    <button class="menu-toggle" on:click={toggleMenu}>☰</button>
    <div class="nav-links" class:is-open={isMenuOpen}>
      {#each pages as page}
        <a href="{base}/{page.slug}" on:click={closeMenu}>{page.title}</a>
      {/each}
      <a href="{base}/blog" on:click={closeMenu}>Blog</a>
    </div>
  </nav>
</header>

<style>
  header {
    margin: 1em 0;
  }

  header a:hover {
    text-decoration: underline;
  }

  nav {
    background-color: var(--header-background);
    padding: 0.75em 2em;
    display: flex;
    justify-content: space-between;
    align-items: center;
    max-width: 1000px;
    margin: 0 auto;
    position: relative;
  }

  .logo {
    font-size: 1.5em;
    font-weight: bold;
    color: var(--tan);
    text-decoration: none;
  }

  .nav-links a {
    color: var(--text-color);
    text-decoration: none;
    display: block;
  }

  .nav-links {
    display: flex;
    gap: 2em;
    font-weight: bold;
  }

  .menu-toggle {
    display: none;
    font-size: 1.5em;
    background: none;
    border: none;
    color: var(--text-color);
    cursor: pointer;
  }

  @media (max-width: 768px) {
    header {
      display: inline-block;
      width: 100%;
      margin: 1px 1px;
    }

    nav {
      padding: 8px 15px;
    }

    .menu-toggle {
      display: block;
      margin-left: auto;
    }

    .nav-links {
      display: none;
      position: absolute;
      top: 100%;
      left: 0;
      right: 0;
      background-color: var(--header-background);
      flex-direction: column;
      gap: 0;
      padding: 1em 0;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
      z-index: 1000;
    }

    .nav-links.is-open {
      display: flex;
    }

    .nav-links a {
      padding: 0.75em 1em;
      text-align: left;
    }

    button {
      padding: 2px 5px;
    }
  }
</style>
