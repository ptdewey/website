import { getCollection } from 'astro:content';
import { fetchStandardPosts } from './standard';

export interface BlogEntry {
  title: string;
  date: Date;
  tags: string[];
  href: string;
  external: boolean;
  source?: string;
}

export async function getAllBlogEntries(): Promise<BlogEntry[]> {
  const local: BlogEntry[] = (await getCollection('blog', ({ data }) => !data.draft)).map(p => ({
    title: p.data.title,
    date: p.data.date,
    tags: p.data.tags ?? [],
    href: `/blog/${p.id}`,
    external: false,
  }));
  const external: BlogEntry[] = (await fetchStandardPosts()).map(p => ({
    title: p.title,
    date: p.date,
    tags: p.tags,
    href: p.href,
    external: true,
    source: p.source,
  }));
  return [...local, ...external].sort((a, b) => b.date.valueOf() - a.date.valueOf());
}
