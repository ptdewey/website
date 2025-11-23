<script lang="ts">
  import FeaturedProjectEntry from "$lib/components/FeaturedProjectEntry.svelte";
  import ProjectEntry from "$lib/components/ProjectEntry.svelte";
  import { featuredProjects, projects } from "./projects";
</script>

<svelte:head>
  <title>Projects</title>
</svelte:head>

<!-- TODO: make separate site with a list of all my tools (like this: https://tools.dhruvs.space/)-->
<div>
  <section>
    <h2>Featured Projects</h2>
    {#each featuredProjects.items as project}
      <FeaturedProjectEntry {project} />
      <hr class="!mt-2 !mb-2" />
    {/each}
    {#each projects as projectCategory}
      <div>
        <h2>{projectCategory.category}</h2>
        <div>
          {#each projectCategory.items as project}
            <ProjectEntry {project} />
          {/each}
        </div>
      </div>
      <hr />
    {/each}
  </section>
</div>
