export const prerender = true;

export async function load() {
  return {
    pages: [{ title: "Projects", slug: "projects" }],
  };
}
