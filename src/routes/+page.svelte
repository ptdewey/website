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
      Hi there! I'm Patrick, and I'm a computer science Masters student at
      Virginia Tech studying software engineering, machine learning, artificial
      intelligence, and data analytics.
    </p>
    <p>
      I love learning new technologies, and I am always working on cool new <a
        href="{base}/projects">projects</a
      >. I am currently working on an intelligent-factory observability platform
      for <a href="https://ccam-va.com/">CCAM</a>, an AI-powered prospecting
      application for
      <a href="https://navagant.com/">Navagant</a> and an experiment on bot detection
      methods for Bluesky.
    </p>
    <p>
      I enjoy optimizing my developer workflow, and I've written Neovim plugins
      and cli applications to try and be blazingly fast in everything I do.
    </p>
    <p>
      Outside of programming, I enjoy reading, listening to music, and making
      coffee and cocktails.
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
      <SocialLink
        title="linkedin"
        link="https://www.linkedin.com/in/patricktdewey/"
      />
      <SocialLink title="bluesky" link="https://bsky.app/profile/pdewey.com" />
      <SocialLink title="rss" link="{base}/rss.xml" />
    </ul>
  </section>
</div>
