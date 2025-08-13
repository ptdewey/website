<script>
  import SocialLink from "$lib/components/SocialLink.svelte";

  import { base } from "$app/paths";
</script>

<svelte:head>
  <title>Patrick Dewey</title>
</svelte:head>

<div>
  <section>
    <h2>Patrick Dewey</h2>
    <p>
      Hello there! I'm Patrick, a backend software engineer who loves learning
      new technologies, building cool <a href="{base}/projects">projects</a>.
    </p>
    <p>
      I also have a <a href="{base}/writing">blog</a> where I
      <em>occasionally</em>
      write articles.
    </p>
  </section>
  <section>
    <ul
      class="font-sm flex flex-col mt-1 mb-0 space-x-0 space-y-2 md:flex-row md:space-x-5 md:space-y-0 list-none pl-0"
    >
      <SocialLink title="github" link="https://github.com/ptdewey" />
      <SocialLink title="bluesky" link="https://bsky.app/profile/pdewey.com" />
      <SocialLink title="rss" link="{base}/rss.xml" />
    </ul>
  </section>
</div>
