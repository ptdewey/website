<script lang="ts">
  let { link, title } = $props();
</script>

<li class="pl-0 prose prose-sm">
  <a class="flex items-center" href={link}>
    <svg
      width="10"
      height="10"
      viewBox="0 0 12 10"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        d="M2.07102 11.3494L0.963068 10.2415L9.2017 1.98864H2.83807L2.85227 0.454545H11.8438V9.46023H10.2955L10.3097 3.09659L2.07102 11.3494Z"
        class="prose prose-sm"
        fill="currentColor"
      >
      </path>
    </svg>
    <p class="ml-2 mt-0 mb-0 prose-a:unset prose prose-sm">{title}</p></a
  >
</li>
