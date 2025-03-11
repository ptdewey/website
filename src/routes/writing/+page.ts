import type { PageLoad } from "./$types";

export const load: PageLoad = async function ({ fetch }) {
  const response = await fetch("/data/writing.json");

  if (!response.ok) {
    throw new Error(`Failed to load posts: ${response.status}`);
  }

  const data = await response.json();

  return { data };
};
