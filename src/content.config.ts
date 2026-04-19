import { defineCollection, z } from 'astro:content';
import { glob } from 'astro/loaders';

const blog = defineCollection({
  loader: glob({ pattern: '**/*.{md,mdx}', base: './src/content/blog' }),
  schema: z.object({
    title: z.string(),
    date: z.coerce.date(),
    authors: z.array(z.string()).optional(),
    categories: z.array(z.string()).optional(),
    tags: z.array(z.string()).optional(),
    type: z.string().optional(),
    description: z.string().optional(),
    draft: z.boolean().optional().default(false),
    bluesky_link: z.string().optional(),
  }),
});

export const collections = { blog };
