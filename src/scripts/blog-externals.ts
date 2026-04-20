import { fetchStandardPosts } from '../lib/standard';

function fmtDate(d: Date): string {
  return d.toLocaleDateString('en-CA', {
    year: 'numeric', month: '2-digit', day: '2-digit', timeZone: 'UTC',
  });
}

function makeEntry(post: {
  title: string;
  date: Date;
  href: string;
  tags: string[];
  source: string;
}): HTMLElement {
  const article = document.createElement('article');
  article.className = 'post-entry';
  article.dataset.entryDate = post.date.toISOString();

  const time = document.createElement('time');
  time.className = 'post-date';
  time.dateTime = post.date.toISOString();
  time.textContent = fmtDate(post.date);

  const meta = document.createElement('div');
  meta.className = 'post-meta';

  const a = document.createElement('a');
  a.className = 'post-title';
  a.href = post.href;
  a.target = '_blank';
  a.rel = 'noopener';
  a.textContent = post.title;
  meta.appendChild(a);

  if (post.tags.length > 0 || post.source) {
    const tags = document.createElement('div');
    tags.className = 'post-tags';
    const src = document.createElement('span');
    src.className = 'source';
    src.textContent = `[${post.source}]`;
    tags.appendChild(src);
    for (const t of post.tags.slice(0, 4)) {
      const tag = document.createElement('span');
      tag.className = 'tag';
      tag.textContent = t;
      tags.appendChild(tag);
    }
    meta.appendChild(tags);
  }

  article.append(time, meta);
  return article;
}

function makeYearSection(year: number): HTMLElement {
  const section = document.createElement('section');
  section.className = 'year-group';
  section.dataset.year = String(year);

  const rule = document.createElement('div');
  rule.className = 'year-rule';
  const y = document.createElement('span');
  y.className = 'year';
  y.textContent = String(year);
  const line = document.createElement('span');
  line.className = 'line';
  line.setAttribute('aria-hidden', 'true');
  const count = document.createElement('span');
  count.className = 'count';
  count.dataset.yearCount = '';
  count.textContent = '0';
  rule.append(y, line, count);

  const entries = document.createElement('div');
  entries.dataset.yearEntries = '';

  section.append(rule, entries);
  return section;
}

function insertIntoYear(section: HTMLElement, entry: HTMLElement, ts: number) {
  const entriesWrap = section.querySelector<HTMLElement>('[data-year-entries]');
  if (!entriesWrap) return;
  const rows = Array.from(entriesWrap.querySelectorAll<HTMLElement>('.post-entry'));
  const before = rows.find(r => {
    const iso = r.dataset.entryDate ?? '';
    return iso && new Date(iso).getTime() < ts;
  });
  if (before) entriesWrap.insertBefore(entry, before);
  else entriesWrap.appendChild(entry);

  const count = section.querySelector<HTMLElement>('[data-year-count]');
  if (count) count.textContent = String(entriesWrap.querySelectorAll('.post-entry').length);
}

function insertYearSection(root: HTMLElement, section: HTMLElement, year: number) {
  const existing = Array.from(root.querySelectorAll<HTMLElement>('.year-group'));
  const before = existing.find(e => Number(e.dataset.year ?? 0) < year);
  if (before) root.insertBefore(section, before);
  else root.appendChild(section);
}

async function hydrate() {
  const root = document.querySelector<HTMLElement>('[data-blog-index]');
  if (!root) return;

  let external;
  try {
    external = await fetchStandardPosts();
  } catch {
    return;
  }
  if (!external.length) return;

  for (const post of external) {
    const year = post.date.getUTCFullYear();
    let section = root.querySelector<HTMLElement>(`.year-group[data-year="${year}"]`);
    if (!section) {
      section = makeYearSection(year);
      insertYearSection(root, section, year);
    }
    insertIntoYear(section, makeEntry(post), post.date.getTime());
  }
}

hydrate();
