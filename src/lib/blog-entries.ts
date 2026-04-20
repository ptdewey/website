import { getCollection } from 'astro:content';

export interface BlogEntry {
  title: string;
  date: Date;
  tags: string[];
  href: string;
  external: boolean;
  source?: string;
}

export async function getLocalBlogEntries(): Promise<BlogEntry[]> {
  const local: BlogEntry[] = (await getCollection('blog', ({ data }) => !data.draft)).map(p => ({
    title: p.data.title,
    date: p.data.date,
    tags: p.data.tags ?? [],
    href: `/blog/${p.id}`,
    external: false,
  }));
  return local.sort((a, b) => b.date.valueOf() - a.date.valueOf());
}
