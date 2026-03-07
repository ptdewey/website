const ATProto = (() => {
  const HANDLE = 'pdewey.com'
  const recordCache = new Map()

  const xrpc = (pds, method, params) =>
    fetch(`${pds}/xrpc/${method}?${new URLSearchParams(params)}`).then(r => r.json())

  const resolveHandle = (pds, handle) =>
    xrpc(pds, 'com.atproto.identity.resolveHandle', {handle: handle || HANDLE}).then(r => r.did)

  const getRecord = (pds, did, collection, atURI) => {
    if (!atURI) return null
    if (recordCache.has(atURI)) return recordCache.get(atURI)
    const p = xrpc(pds, 'com.atproto.repo.getRecord', {
      repo: did, collection, rkey: atURI.split('/').pop()
    }).catch(() => null)
    recordCache.set(atURI, p)
    return p
  }

  const listRecords = (pds, did, collection, limit) =>
    xrpc(pds, 'com.atproto.repo.listRecords', {repo: did, collection, limit})
      .then(r => r.records)

  return { HANDLE, xrpc, resolveHandle, getRecord, listRecords }
})()
