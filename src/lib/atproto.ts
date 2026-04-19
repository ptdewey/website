export const PDS = 'https://arabica.systems';
export const DID = 'did:plc:hm5f3dnm6jdhrc55qp2npdja';
export const HANDLE = 'pdewey.com';

const BUILD_TIMEOUT_MS = 3000;

// In Node (build) we want a hard timeout so a slow PDS never hangs the build.
// In the browser we let fetch run to its own default. Callers pass AbortSignal
// when they want it.
function buildSignal(): AbortSignal | undefined {
  return typeof window === 'undefined' ? AbortSignal.timeout(BUILD_TIMEOUT_MS) : undefined;
}

export async function xrpc<T = unknown>(
  method: string,
  params: Record<string, string | number>,
  opts: { pds?: string; signal?: AbortSignal } = {},
): Promise<T> {
  const pds = opts.pds ?? PDS;
  const qs = new URLSearchParams(
    Object.entries(params).map(([k, v]) => [k, String(v)]),
  );
  const r = await fetch(`${pds}/xrpc/${method}?${qs}`, {
    signal: opts.signal ?? buildSignal(),
  });
  if (!r.ok) throw new Error(`xrpc ${method}: ${r.status}`);
  return (await r.json()) as T;
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

// --- build-time domain helpers (null on failure so a flaky PDS doesn't
// break the static build) ---

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

type Ref = string | { uri?: string; $link?: string };
function refToUri(ref: Ref | undefined | null): string | undefined {
  if (!ref) return undefined;
  if (typeof ref === 'string') return ref;
  return ref.uri ?? ref.$link ?? undefined;
}

interface BrewRecord {
  beanRef?: Ref;
  method?: string;
  createdAt: string;
}

interface BeanRecord {
  name?: string;
}

export async function fetchRecentPlay(): Promise<RecentPlay | null> {
  try {
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
  } catch {
    return null;
  }
}

export async function fetchRecentBrew(): Promise<RecentBrew | null> {
  try {
    const records = await listRecords<BrewRecord>('social.arabica.alpha.brew', 5);
    if (!records.length) return null;
    const sorted = [...records].sort(
      (a, b) => new Date(b.value.createdAt).getTime() - new Date(a.value.createdAt).getTime(),
    );
    const top = sorted[0].value;
    const bean = await getRecord<BeanRecord>('social.arabica.alpha.bean', refToUri(top.beanRef));
    return {
      beanName: bean?.name ?? 'unknown bean',
      method: top.method ?? '',
      createdAt: new Date(top.createdAt),
    };
  } catch {
    return null;
  }
}
