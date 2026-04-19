import { getRecord, listRecords } from "./atproto";

// Publications to skip in the blog index — e.g. a standard.site mirror of this
// site itself, which would otherwise show up once per local post. Accepts
// publication AT URIs (at://did/site.standard.publication/rkey) or host
// substrings matched against the publication's `url` (e.g. 'pdewey.com').
const EXCLUDED_PUBLICATIONS: string[] = ["pdewey.com"];

interface PublicationRecord {
  name?: string;
  url?: string;
}

interface DocumentRecord {
  title?: string;
  description?: string;
  publishedAt?: string;
  path?: string;
  tags?: string[];
  site?: string;
}

export interface ExternalPost {
  title: string;
  description?: string;
  date: Date;
  tags: string[];
  href: string;
  source: string;
}

function parseAtUri(uri: string): { collection: string; rkey: string } | null {
  const m = uri.match(/^at:\/\/[^/]+\/([^/]+)\/([^/]+)$/);
  return m ? { collection: m[1], rkey: m[2] } : null;
}

export async function fetchStandardPosts(limit = 50): Promise<ExternalPost[]> {
  let records: Array<{ uri: string; value: DocumentRecord }> = [];
  try {
    records = await listRecords<DocumentRecord>("site.standard.document", limit);
  } catch {
    return [];
  }
  if (!records.length) return [];

  const pubCache = new Map<string, PublicationRecord | null>();
  const resolvePub = async (uri: string) => {
    if (pubCache.has(uri)) return pubCache.get(uri)!;
    const parsed = parseAtUri(uri);
    const pub = parsed
      ? await getRecord<PublicationRecord>(parsed.collection, parsed.rkey)
      : null;
    pubCache.set(uri, pub);
    return pub;
  };

  const posts: ExternalPost[] = [];
  for (const { value } of records) {
    if (!value.title || !value.publishedAt || !value.site || !value.path) continue;
    if (EXCLUDED_PUBLICATIONS.includes(value.site)) continue;
    const pub = await resolvePub(value.site);
    if (!pub?.url) continue;
    if (
      EXCLUDED_PUBLICATIONS.some(
        (e) => !e.startsWith("at://") && pub.url!.includes(e),
      )
    )
      continue;
    const base = pub.url.replace(/\/$/, "");
    const path = value.path.startsWith("/") ? value.path : `/${value.path}`;
    posts.push({
      title: value.title,
      description: value.description,
      date: new Date(value.publishedAt),
      tags: value.tags ?? [],
      href: `${base}${path}`,
      source: pub.name ?? "standard",
    });
  }
  return posts;
}
