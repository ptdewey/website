<script lang="ts">
  let { project } = $props();
</script>

<div>
  <h3>
    {#if project.link}
      <a href={project.link} target="_blank" rel="noopener noreferrer">
        {project.title}
      </a>
    {:else}
      {project.title}
    {/if}
  </h3>
  {#if project.image}
    <img
      src={project.image}
      alt="{project.title} screenshot"
      class="rounded-sm"
    />
  {/if}

  <p class="mb-0.5">{project.description}</p>
  {#if project.time}
    <p class="italic">{project.time}</p>
  {/if}
</div>
