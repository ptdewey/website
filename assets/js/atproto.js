export const PDS = "https://arabica.systems";
export const DID = "did:plc:hm5f3dnm6jdhrc55qp2npdja";
export const HANDLE = "pdewey.com";

export async function xrpc(method, params, opts = {}) {
  const pds = opts.pds ?? PDS;
  const qs = new URLSearchParams(
    Object.entries(params).map(([k, v]) => [k, String(v)]),
  );
  const r = await fetch(`${pds}/xrpc/${method}?${qs}`, {
    signal: opts.signal,
    cache: "default",
  });
  if (!r.ok) throw new Error(`xrpc ${method}: ${r.status}`);
  return await r.json();
}

export function refToUri(ref) {
  if (!ref) return undefined;
  if (typeof ref === "string") return ref;
  return ref.uri ?? ref.$link ?? undefined;
}

const recordCache = new Map();

export function getRecord(collection, rkeyOrUri, opts = {}) {
  if (!rkeyOrUri) return Promise.resolve(null);
  const rkey = rkeyOrUri.includes("/") ? rkeyOrUri.split("/").pop() : rkeyOrUri;
  const repo = opts.repo ?? DID;
  const key = `${opts.pds ?? PDS}|${repo}|${collection}|${rkey}`;
  const cached = recordCache.get(key);
  if (cached) return cached;
  const p = xrpc("com.atproto.repo.getRecord", { repo, collection, rkey }, opts)
    .then((r) => {
      if (r && typeof r === "object" && "value" in r) return r.value ?? null;
      return r ?? null;
    })
    .catch(() => null);
  recordCache.set(key, p);
  return p;
}

export async function listRecords(collection, limit = 5, opts = {}) {
  const repo = opts.repo ?? DID;
  const res = await xrpc(
    "com.atproto.repo.listRecords",
    { repo, collection, limit },
    opts,
  );
  return res.records;
}

export async function resolveHandle(handle = HANDLE, opts = {}) {
  const r = await xrpc("com.atproto.identity.resolveHandle", { handle }, opts);
  return r.did;
}
