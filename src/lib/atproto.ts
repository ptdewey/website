export const PDS = 'https://arabica.systems';
export const DID = 'did:plc:hm5f3dnm6jdhrc55qp2npdja';
export const HANDLE = 'pdewey.com';

export async function xrpc<T = unknown>(
  method: string,
  params: Record<string, string | number>,
  opts: { pds?: string; signal?: AbortSignal } = {},
): Promise<T> {
  const pds = opts.pds ?? PDS;
  const qs = new URLSearchParams(
    Object.entries(params).map(([k, v]) => [k, String(v)]),
  );
  // cache: 'default' lets the browser honour whatever Cache-Control the
  // PDS returns. All callers are client-side; build-time timeout logic
  // is gone now that nothing runs in Node.
  const r = await fetch(`${pds}/xrpc/${method}?${qs}`, {
    signal: opts.signal,
    cache: 'default',
  });
  if (!r.ok) throw new Error(`xrpc ${method}: ${r.status}`);
  return (await r.json()) as T;
}

export type Ref = string | { uri?: string; $link?: string } | undefined | null;

export function refToUri(ref: Ref): string | undefined {
  if (!ref) return undefined;
  if (typeof ref === 'string') return ref;
  return ref.uri ?? ref.$link ?? undefined;
}

const recordCache = new Map<string, Promise<unknown>>();

export function getRecord<T = unknown>(
  collection: string,
  rkeyOrUri: string | undefined | null,
  opts: { repo?: string; pds?: string } = {},
): Promise<T | null> {
  if (!rkeyOrUri) return Promise.resolve(null);
  const rkey = rkeyOrUri.includes('/') ? rkeyOrUri.split('/').pop()! : rkeyOrUri;
  const repo = opts.repo ?? DID;
  const key = `${opts.pds ?? PDS}|${repo}|${collection}|${rkey}`;
  const cached = recordCache.get(key) as Promise<T | null> | undefined;
  if (cached) return cached;
  const p = xrpc<unknown>('com.atproto.repo.getRecord', { repo, collection, rkey }, opts)
    .then((r): T | null => {
      if (r && typeof r === 'object' && 'value' in r) {
        const v = (r as { value: unknown }).value;
        if (v) return v as T;
      }
      return (r as T) ?? null;
    })
    .catch(() => null);
  recordCache.set(key, p);
  return p;
}

export async function listRecords<T = unknown>(
  collection: string,
  limit = 5,
  opts: { repo?: string; pds?: string } = {},
): Promise<Array<{ uri: string; value: T }>> {
  const repo = opts.repo ?? DID;
  const res = await xrpc<{ records: Array<{ uri: string; value: T }> }>(
    'com.atproto.repo.listRecords',
    { repo, collection, limit },
    opts,
  );
  return res.records;
}

export async function resolveHandle(handle: string = HANDLE, opts: { pds?: string } = {}): Promise<string> {
  const r = await xrpc<{ did: string }>('com.atproto.identity.resolveHandle', { handle }, opts);
  return r.did;
}

// --- domain getters used by the homepage "recent" rows. Both resolve
// to null on any failure so a flaky PDS just leaves the placeholder
// state rather than crashing the script. ---

export interface LatestPlay {
  trackName: string;
  artists: string;
  playedTime: Date;
  uri: string;
}

export interface LatestBrew {
  beanName: string;
  method: string;
  createdAt: Date;
  rkey: string;
}

interface PlayRecord {
  trackName: string;
  playedTime: string;
  artists?: Array<{ artistName: string }>;
}

interface BrewRecord {
  beanRef?: Ref;
  method?: string;
  createdAt: string;
}

interface BeanRecord {
  name?: string;
}

export async function getLatestPlay(): Promise<LatestPlay | null> {
  try {
    const records = await listRecords<PlayRecord>('fm.teal.alpha.feed.play', 10);
    if (!records.length) return null;
    const top = [...records].sort(
      (a, b) => new Date(b.value.playedTime).getTime() - new Date(a.value.playedTime).getTime(),
    )[0];
    return {
      trackName: top.value.trackName,
      artists: top.value.artists?.map(a => a.artistName).join(', ') ?? '',
      playedTime: new Date(top.value.playedTime),
      uri: top.uri,
    };
  } catch {
    return null;
  }
}

export async function getLatestBrew(): Promise<LatestBrew | null> {
  try {
    const records = await listRecords<BrewRecord>('social.arabica.alpha.brew', 5);
    if (!records.length) return null;
    const top = [...records].sort(
      (a, b) => new Date(b.value.createdAt).getTime() - new Date(a.value.createdAt).getTime(),
    )[0];
    const bean = await getRecord<BeanRecord>('social.arabica.alpha.bean', refToUri(top.value.beanRef));
    return {
      beanName: bean?.name ?? 'unknown bean',
      method: top.value.method ?? '',
      createdAt: new Date(top.value.createdAt),
      rkey: top.uri.split('/').pop() ?? '',
    };
  } catch {
    return null;
  }
}
