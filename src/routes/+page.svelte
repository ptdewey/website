<script>
  import { base } from "$app/paths";
</script>

<svelte:head>
  <title>Patrick Dewey</title>
</svelte:head>

<div>
  <section class="welcome">
    <h1>Patrick Dewey</h1>
    <p>
      Hi there! I'm Patrick, and I'm a computer science Masters student at
      Virginia Tech studying artificial intelligence, machine learning, and data
      analytics.
    </p>
    <p>
      I love learning new technologies, and I am always working on cool new <a
        href="{base}/projects">projects</a
      >. I am currently working on an intelligent-factory observability platform
      for <a href="https://ccam-va.com/">CCAM</a>, an AI-powered prospecting
      application for
      <a href="https://navagant.com/">Navagant</a> and a Bluesky bot detection system.
    </p>
    <p>
      I also enjoy optimizing my developer workflow, and I've written vim
      plugins and cli/tui applications to try and be blazingly fast in
      everything I do.
    </p>
    <p>
      I also have a <a href="{base}/blog">blog</a> where I
      <em>occasionally</em>
      write articles. (<a href="{base}/rss.xml">RSS</a>)
    </p>
  </section>
  <section class="contact">
    <p>
      You find me on <a href="https://github.com/ptdewey">GitHub</a>,
      <a href="https://www.linkedin.com/in/patricktdewey/">LinkedIn</a>, and
      <a href="https://bsky.app/profile/pdewey.com">Bluesky</a>
    </p>
  </section>
</div>

<style>
  .welcome {
    max-width: 800px;
  }
</style>
