<script lang="ts">
  import type { Project } from "$lib/types";

  let { project: p } = $props();
  let project: Project = p;
</script>

<li>
  <div class="mb-0">
    <h4 class="prose:unset prose">
      {#if project.link}
        <a href={project.link} target="_blank" rel="noopener noreferrer">
          {project.title}
        </a>
      {:else}
        {project.title}
      {/if}

      {#if project.languages}
        <span class="ml-0.5 mb-0 prose prose-sm italic"
          >({project.languages})</span
        >
      {/if}
    </h4>
  </div>
  <p class="italic !m-0 p-0 prose-sm">
    {project.shortDescription ? project.shortDescription : project.description}
  </p>
</li>
