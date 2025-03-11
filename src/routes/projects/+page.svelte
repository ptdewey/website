<script>
  import ProjectCard from "$lib/components/ProjectCard.svelte";
  import { projects } from "./projects.js";
  import { featuredProjects } from "./projects.js";
</script>

<svelte:head>
  <title>Projects</title>
</svelte:head>

<!--
 TODO:
- Add navigation jumps for moving between sections
-->
<div>
  <div class="category-section">
    <h1 class="category-title">Featured Projects</h1>
    <div class="projects-grid-full">
      {#each featuredProjects.items as project}
        <ProjectCard {project} />
      {/each}
    </div>
  </div>

  {#each projects as projectCategory}
    <div class="category-section">
      <h1 class="category-title">{projectCategory.category}</h1>
      <div class="projects-grid">
        {#each projectCategory.items as project}
          <ProjectCard {project} />
        {/each}
      </div>
    </div>
  {/each}
</div>

<style>
  .category-section {
    margin-bottom: 2rem;
    max-width: 800px;
  }

  .category-title {
    font-size: 1.8em;
    font-weight: bold;
    margin-bottom: 1rem;
    padding: 1rem;
  }

  .projects-grid-full {
    display: grid;
    grid-template-columns: repeat(1, 1fr);
    gap: 1.5rem;
  }

  .projects-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
  }

  @media (max-width: 768px) {
    .projects-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
