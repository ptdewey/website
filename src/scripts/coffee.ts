import { resolveHandle, listRecords, getRecord } from '../lib/atproto';

const NS = 'social.arabica.alpha';
const PDSLS = 'https://pdsls.dev/at/';
const DOT = ' \u00b7 ';

type Ref = string | { uri?: string; $link?: string } | undefined | null;

function refToUri(ref: Ref): string | undefined {
  if (!ref) return undefined;
  if (typeof ref === 'string') return ref;
  return ref.uri ?? ref.$link ?? undefined;
}

interface Pour { waterAmount?: number }
interface BrewRecord {
  createdAt: string;
  beanRef?: Ref;
  brewerRef?: Ref;
  grinderRef?: Ref;
  method?: string;
  coffeeAmount?: number;
  waterAmount?: number;
  temperature?: number;
  timeSeconds?: number;
  pours?: Pour[];
  grindSize?: string | number;
  tastingNotes?: string;
  rating?: number;
}
interface BeanRecord {
  name?: string;
  origin?: string;
  roastLevel?: string;
  process?: string;
  variety?: string;
  roasterRef?: Ref;
}
interface NamedRecord { name?: string }

interface Brew extends BrewRecord {
  rkey: string;
  brewURL: string;
  beanName: string;
  beanOrigin: string;
  beanRoastLevel: string;
  beanProcess: string;
  beanVariety: string;
  roasterName: string;
  brewerName: string;
  grinderName: string;
}

function formatTime(s: number | undefined): string {
  if (!s) return '';
  const m = Math.floor(s / 60);
  return m > 0 ? `${m}:${String(s % 60).padStart(2, '0')}` : `${s}s`;
}

function setOrRemove(el: Element | null, text: string | undefined | null) {
  if (!el) return;
  if (text) el.textContent = text;
  else el.remove();
}

function renderRating(container: Element, rating: number) {
  const n = Math.min(Math.max(rating, 0), 10);
  container.textContent = `${n}/10`;
}

function renderBrew(tmpl: HTMLTemplateElement, brew: Brew, did: string): DocumentFragment {
  const el = tmpl.content.cloneNode(true) as DocumentFragment;

  const bean = el.querySelector<HTMLAnchorElement>('.brew-bean')!;
  bean.textContent = brew.beanName || 'Unknown bean';
  bean.href = brew.brewURL;

  const timeEl = el.querySelector('time')!;
  timeEl.textContent = new Date(brew.createdAt).toLocaleDateString('en-CA', {
    year: 'numeric', month: '2-digit', day: '2-digit',
  });
  timeEl.setAttribute('datetime', brew.createdAt);

  setOrRemove(el.querySelector('.brew-roaster'), brew.roasterName);
  setOrRemove(
    el.querySelector('.brew-sub'),
    [brew.beanOrigin, brew.beanRoastLevel, brew.beanProcess, brew.beanVariety]
      .filter(Boolean)
      .join(DOT),
  );

  const pours = brew.pours?.length;
  const water = brew.waterAmount || brew.pours?.reduce((s, p) => s + (p.waterAmount ?? 0), 0) || 0;
  const ratio = (brew.coffeeAmount && water)
    ? `${brew.coffeeAmount}g / ${water}g`
    : brew.coffeeAmount ? `${brew.coffeeAmount}g` : water ? `${water}g` : null;
  const timeStr = brew.timeSeconds
    ? `${formatTime(brew.timeSeconds)}${pours ? ` (${pours} pour${pours === 1 ? '' : 's'})` : ''}`
    : null;
  el.querySelector('.brew-meta')!.textContent = [
    brew.method,
    ratio,
    brew.temperature ? `${(brew.temperature / 10).toFixed(1)}\u00b0C` : null,
    timeStr,
  ].filter(Boolean).join(DOT);

  const grinderStr = brew.grinderName
    ? `${brew.grinderName}${brew.grindSize ? ` (${brew.grindSize})` : ''}`
    : (brew.grindSize ? `grind ${brew.grindSize}` : null);
  setOrRemove(el.querySelector('.brew-equipment'), [brew.brewerName, grinderStr].filter(Boolean).join(DOT));

  setOrRemove(el.querySelector('.brew-notes'), brew.tastingNotes);

  const ratingEl = el.querySelector('.brew-rating')!;
  if (brew.rating) renderRating(ratingEl, brew.rating);
  else ratingEl.remove();

  const rkeyEl = el.querySelector('.brew-rkey')!;
  const label = document.createElement('span');
  label.textContent = 'rec  ';
  const a = document.createElement('a');
  a.href = `${PDSLS}${did}/${NS}.brew/${brew.rkey}`;
  a.target = '_blank';
  a.rel = 'noopener';
  a.textContent = brew.rkey;
  rkeyEl.append(label, a);

  return el;
}

async function loadBrews() {
  const container = document.getElementById('brew-list')!;
  const tmpl = document.getElementById('brew-tmpl') as HTMLTemplateElement;
  const setMessage = (text: string) => {
    const p = document.createElement('p');
    p.className = 'log-msg';
    p.textContent = text;
    container.replaceChildren(p);
  };

  try {
    const did = await resolveHandle();
    const records = await listRecords<BrewRecord>(`${NS}.brew`, 100, { repo: did });

    const sorted = records.sort(
      (a, b) => new Date(b.value.createdAt).getTime() - new Date(a.value.createdAt).getTime(),
    );
    const recent = sorted.slice(0, 10);

    const brews: Brew[] = await Promise.all(recent.map(async ({ uri, value }) => {
      const [beanRec, brewerRec, grinderRec] = await Promise.all([
        getRecord<BeanRecord>(`${NS}.bean`, refToUri(value.beanRef), { repo: did }),
        getRecord<NamedRecord>(`${NS}.brewer`, refToUri(value.brewerRef), { repo: did }),
        getRecord<NamedRecord>(`${NS}.grinder`, refToUri(value.grinderRef), { repo: did }),
      ]);
      const roasterRec = await getRecord<NamedRecord>(
        `${NS}.roaster`, refToUri(beanRec?.roasterRef), { repo: did },
      );
      const rkey = uri.split('/').pop() ?? '';
      return {
        ...value,
        rkey,
        brewURL: `https://alpha.arabica.social/brews/${rkey}`,
        beanName: beanRec?.name ?? '',
        beanOrigin: beanRec?.origin ?? '',
        beanRoastLevel: beanRec?.roastLevel ?? '',
        beanProcess: beanRec?.process ?? '',
        beanVariety: beanRec?.variety ?? '',
        roasterName: roasterRec?.name ?? '',
        brewerName: brewerRec?.name ?? '',
        grinderName: grinderRec?.name ?? '',
      };
    }));

    if (brews.length === 0) { setMessage('No brews yet.'); return; }

    container.replaceChildren(...brews.map(b => renderBrew(tmpl, b, did)));
  } catch (e) {
    setMessage('Could not load brews.');
    console.error('Failed to load brews:', e);
  }
}

loadBrews();
