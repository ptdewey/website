---
title: "Adding Comments to My Blog with Bluesky!"
authors: ["Patrick Dewey"]
date: 2024-11-26
categories: ["Software Development"]
tags: ["bluesky", "software-development", "svelte", "npm"]
type: "Blog"
bluesky_link: "at://did:plc:hm5f3dnm6jdhrc55qp2npdja/app.bsky.feed.post/3lbvhkkrg4c2c"
---

### Introduction

Bluesky has been gaining a lot of traction and publicity recently, and there are a few aspects of the platform that have piqued my interest in a way that no other social media platform has done.

What gets me most excited about Bluesky is that both the [platform](https://bsky.app) and the underlying [AT Protocol](https://atproto.com/) are open source[^app] [^atproto]. The open source nature of the platform is great for transparency (something other platforms sorely lack), and also makes it really easy for developers to improve on and build off of the platform. (The AT protocol is also very interesting and will likely enable many exciting projects in the future, but that's a topic for another day)

What I'd like to talk about today is one specific example of a tool that was built off of Bluesky, embeddable blog post comments.

### Embedding Bluesky for Blog Comments

The idea of using Bluesky threads as blog comments was first debuted by [Emily Liu](https://bsky.app/profile/emilyliu.me) (a Bluesky developer) in a [post](https://emilyliu.me/blog/open-network) a couple of days ago. It quickly got people very excited, with many requests for example code.

Fortunately, she published her code in a new [post](https://bsky.app/profile/emilyliu.me/post/3lbqta5lnck2i) the very next day. Unfortunately, her code[^em-gist] was built specifically around NextJS and TailwindCSS, but this code was enough of a starting point for me and (presumably) a number of other developers to start adapting it for other frameworks and applications.

I found [Cory Zue's](https://bsky.app/profile/coryzue.com) implementation[^co] in React pretty early on in its development, but I wasn't able to get it working with my site, as (at the time) the library only worked with React projects. As such, I set out to make my own version of the library in Svelte, since that is the framework I use for my site.

#### Building the Library

To begin, I used the Svelte library creation tool (`npx sv create`) to set up a base project (I chose to use JavaScript and JSDoc for types).

In a similar design to Emily Liu, I created Svelte components for 'Actions', 'Comments', and the 'CommentSection', with the 'CommentSection' being the main exported component to be used in other pages. I wrote the component in a way that allows for either a post URI or a Bluesky handle to be used, with the handle option using the Bluesky API to find the most recent post containing the blog post URL (returning an error message if none are found).

I set up some basic CSS styling for the components, and created some variables to allow styling of certain aspects to be outsourced to library users. This part of the library could likely be improved; I was unsure about which parts of the styling should be handled on the library's end instead of being delegated to users. I hope to address this soon (I'd also very happily accept [feature requests](https://github.com/ptdewey/bluesky-comments-svelte/issues) and [contributions](https://github.com/ptdewey/bluesky-comments-svelte/pulls)).

The last step was publishing the library so that I (and hopefully others) could use the library. I had never published an NPM package before, but the process was actually quite simple, as the library creation tool provided the necessary scripts. After running `npm package` (and making some tweaks to 'package.json'), my library was published to the NPM registry[^npm-bcs]. (The [source code](https://github.com/ptdewey/bluesky-comments-svelte) is also available[^bcs])

### Integrating the Library Into My Site

To get the newly published library integrated with my site, all I had to do was install it:
```sh
npm install bluesky-comments-svelte
```
Import the component into my blog post '+page.svelte' file and set the author prop:
```svelte
<script>
  import { CommentSection } from "bluesky-comments-svelte";
  const author = "pdewey.com";
</script>
```
And then add the component to the page itself:
```svelte
<div>
  <h2>Comments</h2>
  <div class="comment-section">
    <CommentSection {author} opts={{ showCommentsTitle: false }} />
  </div>
</div>
```

With just a few styling modifications (mostly colors), everything looked the way I wanted it.
```svelte
<style>
  :root {
    --comment-border-color: var(--tan);
    --avatar-size: 2rem;
    --font-size-title: 1.5rem;
    --handle-color: var(--green);
    --comment-content-alignment: flex-start;
    --font-size-comment-body: 1rem;
  }

  .comment-section {
    background-color: var(--dark-brown-alt);
    padding: 10px;
    border-radius: 8px;
  }
</style>
```

All that was left to take advantage of it was creating a new post (on the blog and on Bluesky)!

### Closing Thoughts

Bluesky development feels like the wild west right now with many developers creating innovative new tools and integrations for and with the platform, and I'm enjoying being a part of it. I had a lot of fun developing my library and learning about how the Bluesky API works, and I'm quite happy with the result.

If you use Svelte for your website, and want to integrate a comments section powered by Bluesky, give my library a try (and feel free to make suggestions that you think would improve the experience).
If you are not a Svelte user, but you want to add Bluesky comments, I also recommend Cory Zue's `bluesky-comments` library which is available through NPM[^co-npm], and I think it can be used with any framework now as well (although the bundle size is considerably bigger than the Svelte version if that is something you care about).

Overall, this was a great learning experience, and I look forward to seeing all the cool stuff that gets built around Bluesky.

[^em-post]: Emily Liu's original blog post - [https://emilyliu.me/blog/open-network](https://emilyliu.me/blog/open-network)
[^em-gist]: Emily Liu's comment code - [https://gist.github.com/emilyliu7321/19ac4e111588bdc0cb4e411c88d9c79a](https://gist.github.com/emilyliu7321/19ac4e111588bdc0cb4e411c88d9c79a)
[^co]: Cory Zue's comment library - [https://bsky.app/profile/coryzue.com/post/3lbrkypd37224](https://bsky.app/profile/coryzue.com/post/3lbrkypd37224)
[^co-npm]: Cory Zue's library (NPM) - [https://www.npmjs.com/package/bluesky-comments](https://www.npmjs.com/package/bluesky-comments)
[^app]: Bluesky appplication repository - [https://github.com/bluesky-social/social-app](https://github.com/bluesky-social/social-app)
[^atproto]: AT Protocol repository - [https://github.com/bluesky-social/atproto](https://github.com/bluesky-social/atproto)
[^bcs]: My project repository - [https://github.com/ptdewey/bluesky-comments-svelte](https://github.com/ptdewey/bluesky-comments-svelte)
[^npm-bcs]: bluesky-comments-svelte NPM package - [https://www.npmjs.com/package/bluesky-comments-svelte](https://www.npmjs.com/package/bluesky-comments-svelte)
