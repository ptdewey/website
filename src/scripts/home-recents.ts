import { getLatestPlay, getLatestBrew } from '../lib/atproto';
import { fetchStandardPosts } from '../lib/standard';
import { fadeIn, flipAnimate, swap } from '../lib/anim';

function timeAgo(d: Date): string {
  const s = Math.floor((Date.now() - d.getTime()) / 1000);
  if (s < 60) return 'just now';
  const m = Math.floor(s / 60);
  if (m < 60) return `~${m}m ago`;
  const h = Math.floor(m / 60);
  if (h < 24) return `~${h}h ago`;
  const dy = Math.floor(h / 24);
  if (dy === 1) return 'yesterday';
  return `${dy}d ago`;
}

function fmtDate(d: Date): string {
  return d.toLocaleDateString('en-CA', {
    year: 'numeric', month: '2-digit', day: '2-digit', timeZone: 'UTC',
  });
}

function secondaryLine(
  host: HTMLElement,
  parts: Array<{ text: string; className?: string }>,
  timeIso: string,
  timeText: string,
) {
  host.textContent = '';
  for (const part of parts) {
    const span = document.createElement('span');
    if (part.className) span.className = part.className;
    span.textContent = part.text;
    host.appendChild(span);
    const sep = document.createElement('span');
    sep.className = 'sep';
    sep.setAttribute('aria-hidden', 'true');
    sep.textContent = ' · ';
    host.appendChild(sep);
  }
  const t = document.createElement('time');
  t.dateTime = timeIso;
  t.textContent = timeText;
  host.appendChild(t);
}

function fillPlay() {
  const track = document.querySelector<HTMLElement>('[data-play-track]');
  const secondary = document.querySelector<HTMLElement>('[data-play-secondary]');
  if (!track || !secondary) return;

  getLatestPlay().then(play => {
    if (!play) {
      track.textContent = 'no recent play';
      secondary.textContent = '';
      fadeIn(track);
      return;
    }
    track.textContent = play.trackName;
    secondaryLine(
      secondary,
      play.artists ? [{ text: play.artists }] : [],
      play.playedTime.toISOString(),
      timeAgo(play.playedTime),
    );
    fadeIn(track);
    fadeIn(secondary);
  });
}

function fillBrew() {
  const bean = document.querySelector<HTMLElement>('[data-brew-bean]');
  const secondary = document.querySelector<HTMLElement>('[data-brew-secondary]');
  if (!bean || !secondary) return;

  getLatestBrew().then(brew => {
    if (!brew) {
      bean.textContent = 'no recent brew';
      secondary.textContent = '';
      fadeIn(bean);
      return;
    }
    bean.textContent = brew.beanName;
    secondaryLine(
      secondary,
      brew.method ? [{ text: brew.method, className: 'brew-method' }] : [],
      brew.createdAt.toISOString(),
      timeAgo(brew.createdAt),
    );
    fadeIn(bean);
    fadeIn(secondary);
  });
}

function makeBlogEntry(post: {
  title: string;
  date: Date;
  href: string;
  source?: string;
}): HTMLDivElement {
  const row = document.createElement('div');
  row.className = 'entry';
  row.dataset.entryDate = post.date.toISOString();

  const time = document.createElement('time');
  time.dateTime = post.date.toISOString();
  time.textContent = fmtDate(post.date);

  const a = document.createElement('a');
  a.href = post.href;
  a.target = '_blank';
  a.rel = 'noopener';
  a.textContent = post.title;
  if (post.source) {
    const tag = document.createElement('span');
    tag.className = 'source-tag';
    tag.textContent = ` [${post.source}]`;
    a.appendChild(tag);
  }

  row.append(time, a);
  return row;
}

function hydratePosts() {
  const latestRow = document.querySelector<HTMLElement>('[data-latest-post]');
  const blogList = document.querySelector<HTMLElement>('[data-blog-list]');

  fetchStandardPosts().then(external => {
    if (!external.length) return;
    const sorted = [...external].sort((a, b) => b.date.valueOf() - a.date.valueOf());

    // latest-post row: swap in (crossfade) if any external predates current latest.
    if (latestRow) {
      const currentIso = latestRow.dataset.latestDate ?? '';
      const currentTime = currentIso ? new Date(currentIso).getTime() : 0;
      const top = sorted[0];
      if (top.date.getTime() > currentTime) {
        const stack = latestRow.querySelector<HTMLElement>('.value-stack');
        const link = latestRow.querySelector<HTMLAnchorElement>('[data-latest-link]');
        const secondary = latestRow.querySelector<HTMLElement>('[data-latest-secondary]');
        if (stack && link && secondary) {
          swap(stack, () => {
            link.textContent = top.title;
            link.href = top.href;
            link.target = '_blank';
            link.rel = 'noopener';
            secondary.textContent = '';
            if (top.source) {
              secondary.appendChild(document.createTextNode(`[${top.source}]`));
              const sep = document.createElement('span');
              sep.className = 'sep';
              sep.setAttribute('aria-hidden', 'true');
              sep.textContent = ' · ';
              secondary.appendChild(sep);
            }
            const t = document.createElement('time');
            t.dateTime = top.date.toISOString();
            t.textContent = fmtDate(top.date);
            secondary.appendChild(t);
            latestRow.dataset.latestDate = top.date.toISOString();
          });
        }
      }
    }

    // blog 3-list: insert externals in date order, trim to 3 entries.
    // Displaced siblings slide to their new positions; the new entry fades in.
    if (blogList) {
      for (const post of sorted) {
        const children = Array.from(blogList.querySelectorAll<HTMLElement>('.entry'));
        const target = children.find(c => {
          const iso = c.dataset.entryDate ?? '';
          return iso && new Date(iso).getTime() < post.date.getTime();
        });
        if (!target && children.length >= 3) continue;

        const entry = makeBlogEntry(post);
        flipAnimate(children, () => {
          if (target) blogList.insertBefore(entry, target);
          else blogList.appendChild(entry);
          const all = Array.from(blogList.querySelectorAll<HTMLElement>('.entry'));
          for (const extra of all.slice(3)) extra.remove();
        });
        fadeIn(entry);
      }
    }
  }).catch(() => {
    // silently ignore — SSR'd local posts remain visible.
  });
}

fillPlay();
fillBrew();
hydratePosts();
