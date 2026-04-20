import { xrpc, listRecords, DID } from '../lib/atproto';

const NS = 'fm.teal.alpha';
const PDSLS = 'https://pdsls.dev/at/';
const BAR_WIDTH = 24;

interface Artist { artistName: string }
interface PlayItem {
  trackName: string;
  artists?: Artist[];
  releaseName?: string;
  duration?: number;
}
interface PlayRecord {
  trackName: string;
  artists?: Artist[];
  releaseName?: string;
  playedTime: string;
  originUrl?: string;
}
interface StatusRecord {
  item?: PlayItem;
  time?: string | number;
  expiry?: string | number;
}

function fmtMs(s: number): string {
  const m = Math.floor(s / 60);
  return `${m}:${String(Math.floor(s % 60)).padStart(2, '0')}`;
}

function fmtStamp(iso: string): string {
  const d = new Date(iso);
  const now = new Date();
  const sameDay =
    d.getFullYear() === now.getFullYear() &&
    d.getMonth() === now.getMonth() &&
    d.getDate() === now.getDate();
  const time = d.toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' });
  if (sameDay) return time;
  const date = d.toLocaleDateString('en-CA', { month: '2-digit', day: '2-digit' });
  return `${date} ${time}`;
}

function renderPlay(tmpl: HTMLTemplateElement, uri: string, play: PlayRecord): DocumentFragment {
  const el = tmpl.content.cloneNode(true) as DocumentFragment;

  const link = el.querySelector<HTMLAnchorElement>('.play-track')!;
  link.textContent = play.trackName;
  if (play.originUrl) link.href = play.originUrl;
  else link.removeAttribute('target');

  const timeEl = el.querySelector('time')!;
  timeEl.textContent = fmtStamp(play.playedTime);
  timeEl.setAttribute('datetime', play.playedTime);

  el.querySelector('.play-artist')!.textContent =
    play.artists?.map(a => a.artistName).join(', ') ?? '';

  const releaseEl = el.querySelector('.play-release')!;
  if (play.releaseName) releaseEl.textContent = play.releaseName;
  else releaseEl.remove();

  const rkeyEl = el.querySelector('.play-rkey')!;
  const rkey = uri.split('/').pop() ?? '';
  const label = document.createElement('span');
  label.textContent = 'rec  ';
  const a = document.createElement('a');
  a.href = `${PDSLS}${DID}/${NS}.feed.play/${rkey}`;
  a.target = '_blank';
  a.rel = 'noopener';
  a.textContent = rkey;
  rkeyEl.append(label, a);

  return el;
}

let npTimer: ReturnType<typeof setInterval> | null = null;

async function loadStatus() {
  if (npTimer) { clearInterval(npTimer); npTimer = null; }

  try {
    const rec = await xrpc<unknown>('com.atproto.repo.getRecord', {
      repo: DID,
      collection: `${NS}.actor.status`,
      rkey: 'self',
    }).catch(() => null);
    const status: StatusRecord | undefined = rec && typeof rec === 'object'
      ? ('value' in rec && (rec as { value: unknown }).value
          ? (rec as { value: StatusRecord }).value
          : (rec as StatusRecord))
      : undefined;
    const container = document.getElementById('now-playing')!;
    if (!status?.item) { container.classList.add('hidden'); return; }

    const startTime = Number(status.time);
    const expiryTime = Number(status.expiry);
    const now = Math.floor(Date.now() / 1000);
    if (expiryTime && expiryTime < now) { container.classList.add('hidden'); return; }

    const total = status.item.duration || null;
    document.getElementById('np-track')!.textContent = status.item.trackName;
    document.getElementById('np-artist')!.textContent =
      status.item.artists?.map(a => a.artistName).join(', ') ?? '';
    const releaseEl = document.getElementById('np-release')!;
    if (status.item.releaseName) {
      releaseEl.textContent = status.item.releaseName;
      releaseEl.classList.remove('hidden');
    } else {
      releaseEl.textContent = '';
    }

    const timeEl = document.getElementById('np-time')!;
    const filledEl = document.getElementById('np-bar-filled')!;
    const emptyEl = document.getElementById('np-bar-empty')!;

    const updateTime = () => {
      const elapsed = Math.max(0, Math.floor(Date.now() / 1000) - startTime);
      if (total) {
        const clamped = Math.min(elapsed, total);
        const ratio = clamped / total;
        const fill = Math.max(0, Math.min(BAR_WIDTH, Math.round(ratio * BAR_WIDTH)));
        filledEl.textContent = '\u2588'.repeat(fill);
        emptyEl.textContent = '\u2591'.repeat(BAR_WIDTH - fill);
        timeEl.textContent = `${fmtMs(clamped)} / ${fmtMs(total)}`;
        if (clamped >= total) {
          if (npTimer) clearInterval(npTimer);
          npTimer = null;
          loadStatus();
          loadPlays();
        }
      } else {
        filledEl.textContent = '';
        emptyEl.textContent = '\u2591'.repeat(BAR_WIDTH);
        timeEl.textContent = fmtMs(elapsed);
      }
    };
    updateTime();
    npTimer = setInterval(updateTime, 1000);

    container.classList.remove('hidden');
  } catch (e) {
    console.error('Failed to load status:', e);
  }
}

async function loadPlays() {
  const container = document.getElementById('play-list')!;
  const tmpl = document.getElementById('play-tmpl') as HTMLTemplateElement;
  const setMessage = (text: string) => {
    const p = document.createElement('p');
    p.className = 'log-msg';
    p.textContent = text;
    container.replaceChildren(p);
  };

  try {
    const records = await listRecords<PlayRecord>(`${NS}.feed.play`, 15);
    const sorted = records.sort(
      (a, b) => new Date(b.value.playedTime).getTime() - new Date(a.value.playedTime).getTime(),
    );
    const recent = sorted.slice(0, 10);

    if (recent.length === 0) { setMessage('No plays yet.'); return; }

    container.replaceChildren(...recent.map(r => renderPlay(tmpl, r.uri, r.value)));
  } catch (e) {
    setMessage('Could not load plays.');
    console.error('Failed to load plays:', e);
  }
}

loadStatus();
loadPlays();
