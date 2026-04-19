const PDS = 'https://arabica.systems';
const DID = 'did:plc:hm5f3dnm6jdhrc55qp2npdja';
const TIMEOUT_MS = 3000;

async function xrpc<T>(method: string, params: Record<string, string | number>): Promise<T | null> {
  const qs = new URLSearchParams(
    Object.entries(params).map(([k, v]) => [k, String(v)]),
  );
  const controller = new AbortController();
  const t = setTimeout(() => controller.abort(), TIMEOUT_MS);
  try {
    const r = await fetch(`${PDS}/xrpc/${method}?${qs}`, { signal: controller.signal });
    if (!r.ok) return null;
    return (await r.json()) as T;
  } catch {
    return null;
  } finally {
    clearTimeout(t);
  }
}

async function getRecord<T>(collection: string, rkey: string): Promise<T | null> {
  const res = await xrpc<{ value: T }>('com.atproto.repo.getRecord', {
    repo: DID,
    collection,
    rkey,
  });
  return res?.value ?? null;
}

async function listRecords<T>(collection: string, limit = 5): Promise<Array<{ uri: string; value: T }>> {
  const res = await xrpc<{ records: Array<{ uri: string; value: T }> }>(
    'com.atproto.repo.listRecords',
    { repo: DID, collection, limit },
  );
  return res?.records ?? [];
}

export interface RecentPlay {
  trackName: string;
  artists: string;
  playedTime: Date;
}

export interface RecentBrew {
  beanName: string;
  method: string;
  createdAt: Date;
}

interface PlayRecord {
  trackName: string;
  playedTime: string;
  artists?: Array<{ artistName: string }>;
}

interface BrewRecord {
  beanRef?: string;
  method?: string;
  createdAt: string;
}

interface BeanRecord {
  name?: string;
}

export async function fetchRecentPlay(): Promise<RecentPlay | null> {
  const records = await listRecords<PlayRecord>('fm.teal.alpha.feed.play', 20);
  if (!records.length) return null;
  const sorted = [...records].sort(
    (a, b) => new Date(b.value.playedTime).getTime() - new Date(a.value.playedTime).getTime(),
  );
  const top = sorted[0].value;
  return {
    trackName: top.trackName,
    artists: top.artists?.map(a => a.artistName).join(', ') ?? '',
    playedTime: new Date(top.playedTime),
  };
}

export async function fetchRecentBrew(): Promise<RecentBrew | null> {
  const records = await listRecords<BrewRecord>('social.arabica.alpha.brew', 5);
  if (!records.length) return null;
  const sorted = [...records].sort(
    (a, b) => new Date(b.value.createdAt).getTime() - new Date(a.value.createdAt).getTime(),
  );
  const top = sorted[0].value;
  const rkey = top.beanRef?.split('/').pop();
  const bean = rkey ? await getRecord<BeanRecord>('social.arabica.alpha.bean', rkey) : null;
  return {
    beanName: bean?.name ?? 'unknown bean',
    method: top.method ?? '',
    createdAt: new Date(top.createdAt),
  };
}
