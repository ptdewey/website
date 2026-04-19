// @ts-check
import { defineConfig } from 'astro/config';
import mdx from '@astrojs/mdx';
import tailwindcss from '@tailwindcss/vite';

/**
 * Inline rehype plugin that adds loading="lazy" and decoding="async"
 * to every <img> in rendered markdown. Below-the-fold post images
 * shouldn't block the first paint.
 */
const rehypeLazyImages = () => (tree) => {
  const walk = (node) => {
    if (node.type === 'element' && node.tagName === 'img') {
      node.properties = node.properties ?? {};
      node.properties.loading = 'lazy';
      node.properties.decoding = 'async';
    }
    if (node.children) node.children.forEach(walk);
  };
  walk(tree);
};

// https://astro.build/config
export default defineConfig({
  site: 'https://pdewey.com',
  integrations: [mdx()],
  vite: {
    plugins: [tailwindcss()],
  },
  markdown: {
    shikiConfig: {
      theme: 'github-dark',
    },
    rehypePlugins: [rehypeLazyImages],
  },
});
