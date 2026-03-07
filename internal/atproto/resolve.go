package atproto

import (
	"context"
	"fmt"

	"github.com/bluesky-social/indigo/atproto/identity"
	"github.com/bluesky-social/indigo/atproto/syntax"
)

// dir is the package-level identity directory used for handle and DID resolution.
// DefaultDirectory returns a CacheDirectory with DNS+HTTPS handle resolution,
// LRU caching, and support for both did:plc and did:web.
var dir = identity.DefaultDirectory()

// ResolveHandle resolves an ATProto handle to a DID string.
func ResolveHandle(ctx context.Context, handle string) (string, error) {
	h, err := syntax.ParseHandle(handle)
	if err != nil {
		return "", fmt.Errorf("invalid handle %q: %w", handle, err)
	}
	ident, err := dir.LookupHandle(ctx, h)
	if err != nil {
		return "", fmt.Errorf("resolving handle %q: %w", handle, err)
	}
	return ident.DID.String(), nil
}

// ResolvePDS resolves a DID string to its PDS service endpoint URL.
func ResolvePDS(ctx context.Context, did string) (string, error) {
	d, err := syntax.ParseDID(did)
	if err != nil {
		return "", fmt.Errorf("invalid DID %q: %w", did, err)
	}
	ident, err := dir.LookupDID(ctx, d)
	if err != nil {
		return "", fmt.Errorf("resolving DID %q: %w", did, err)
	}
	pds := ident.PDSEndpoint()
	if pds == "" {
		return "", fmt.Errorf("no PDS service found for DID %s", did)
	}
	return pds, nil
}

// ResolveHandleAndPDS resolves a handle to its DID and PDS endpoint in one call.
func ResolveHandleAndPDS(ctx context.Context, handle string) (did, pds string, err error) {
	h, err := syntax.ParseHandle(handle)
	if err != nil {
		return "", "", fmt.Errorf("invalid handle %q: %w", handle, err)
	}
	ident, err := dir.LookupHandle(ctx, h)
	if err != nil {
		return "", "", fmt.Errorf("resolving handle %q: %w", handle, err)
	}
	pds = ident.PDSEndpoint()
	if pds == "" {
		return "", "", fmt.Errorf("no PDS service found for handle %s", handle)
	}
	return ident.DID.String(), pds, nil
}
