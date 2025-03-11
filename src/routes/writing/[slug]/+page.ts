import type { Writing } from "$lib/types";
import type { PageLoad } from "./$types";

export const load: PageLoad = async function ({ params, fetch }: any) {
  const { slug } = params;

  const response = await fetch("/data/writing.json");

  if (!response.ok) {
    throw new Error(`Failed to load posts: ${response.status}`);
  }

  const writings: Writing[] = await response.json();

  const data = writings.find((w: Writing) => w.metadata.slug === slug);

  if (!data) {
    return {
      status: 404,
      error: new Error("Post not found"),
    };
  }

  return data;
};
