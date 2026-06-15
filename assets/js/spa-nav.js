const importedModules = new Set();

async function hydratePage(container, baseUrl = document.baseURI) {
  const scripts = [
    ...container.querySelectorAll('script[type="module"][src]'),
  ];

  await Promise.all(
    scripts.map((script) => {
      const src = new URL(script.getAttribute("src"), baseUrl).href;

      if (importedModules.has(src)) {
        return Promise.resolve();
      }

      importedModules.add(src);
      return import(src).catch((error) => {
        importedModules.delete(src);
        throw error;
      });
    })
  );
}

async function navigateTo(url, push = true) {
  const response = await fetch(url, {
    headers: { "X-Requested-With": "spa-nav" },
  });

  if (!response.ok) {
    location.href = url;
    return;
  }

  const html = await response.text();
  const doc = new DOMParser().parseFromString(html, "text/html");

  const currentPage = document.querySelector("#page");
  const nextPage = doc.querySelector("#page");

  if (!currentPage || !nextPage) {
    location.href = url;
    return;
  }

  document.title = doc.title;
  currentPage.innerHTML = nextPage.innerHTML;
  try {
    await hydratePage(currentPage, url);
  } catch (error) {
    console.warn("spa-nav: failed to hydrate page scripts", error);
  }

  if (push) {
    history.pushState(null, "", url);
  }

  window.scrollTo(0, 0);
}

document.addEventListener("click", (event) => {
  const link = event.target.closest("a");
  if (!link) return;

  const url = new URL(link.href);

  if (url.origin !== location.origin) return;
  if (link.target || link.hasAttribute("download")) return;
  if (url.hash && url.pathname === location.pathname) return;
  if (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;

  event.preventDefault();
  navigateTo(url.href);
});

window.addEventListener("popstate", () => {
  navigateTo(location.href, false);
});
