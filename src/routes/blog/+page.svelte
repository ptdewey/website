<script>
  import { base } from "$app/paths";
  export let data;
</script>

<svelte:head>
  <title>Blog</title>
</svelte:head>

<div class="blog-list-container">
  <h1>Blog</h1>
  <ul>
    {#each data.posts as post}
      <li>
        {post.metadata.date} -
        <a href="{base}/blog/{post.metadata.slug}">
          {post.metadata.title}
        </a>
      </li>
    {/each}
  </ul>
</div>

<style>
  .blog-list-container {
    max-width: 1200px;
  }

  li {
    list-style: none;
  }
</style>
