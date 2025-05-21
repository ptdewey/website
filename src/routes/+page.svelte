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
      Hello there! I'm Patrick, a backend software engineer at Koalafi and a
      recent graduate of Virginia Tech with a Master's degree in Computer
      Science.
    </p>
    <p>
      I love learning new technologies, and I am always working on cool new <a
        href="{base}/projects">projects</a
      >. I am currently working on a
      <a href="https://github.com/ptdewey/plantuml-lsp"
        >language server for PlantUML</a
      >
      and am in the process of improving my
      <a href="https://github.com/ptdewey/pendulum-nvim"
        >time tracking plugin for Neovim</a
      >.
    </p>
    <p>
      I enjoy optimizing my developer workflow, and I've written Neovim plugins
      and cli applications to be blazingly fast in everything I do.
    </p>
    <p>
      Outside of programming, I enjoy reading, listening to music, playing
      Minecraft, and making coffee and cocktails.
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
