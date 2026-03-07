package atproto

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bluesky-social/indigo/atproto/atclient"
	"github.com/bluesky-social/indigo/atproto/syntax"
)

// PublicationPageInfo holds display info for a single document fetched from the PDS.
type PublicationPageInfo struct {
	Title       string
	Description string
	Link        string    // full public URL: publication.url + record path
	Date        time.Time // parsed from publishedAt
	ATURI       string    // AT-URI of the record
}

// listRecordsResponse is the shape returned by com.atproto.repo.listRecords.
// We use our own struct (rather than the generated comatproto.RepoListRecords_Output)
// because site.standard.document is not registered in indigo's lexicon type map,
// which would cause LexiconTypeDecoder.UnmarshalJSON to return an error.
type listRecordsResponse struct {
	Records []recordEntry `json:"records"`
	Cursor  string        `json:"cursor,omitempty"`
}

type recordEntry struct {
	URI   string          `json:"uri"`
	CID   string          `json:"cid"`
	Value json.RawMessage `json:"value"`
}

type documentValue struct {
	Type        string `json:"$type"`
	Title       string `json:"title"`
	PublishedAt string `json:"publishedAt"`
	Path        string `json:"path"`
	Description string `json:"description,omitempty"`
	Site        string `json:"site,omitempty"`
}

// FetchPublicationDocuments fetches all site.standard.document records for the
// given DID from the PDS and returns them sorted by date descending.
// pdsURL is the base URL of the PDS (e.g. "https://bsky.social").
// pubURL is the publication's public URL (e.g. "https://mypub.leaflet.pub"),
// used to construct the full link from each record's path.
func FetchPublicationDocuments(ctx context.Context, pdsURL, did, pubURL string) ([]PublicationPageInfo, error) {
	c := atclient.NewAPIClient(pdsURL)

	nsid, err := syntax.ParseNSID("com.atproto.repo.listRecords")
	if err != nil {
		return nil, fmt.Errorf("parsing NSID: %w", err)
	}

	params := map[string]any{
		"repo":       did,
		"collection": "site.standard.document",
		"limit":      int64(100),
	}

	var result listRecordsResponse
	if err := c.Get(ctx, nsid, params, &result); err != nil {
		return nil, fmt.Errorf("fetching publication documents: %w", err)
	}

	pubBase := strings.TrimRight(pubURL, "/")
	docs := make([]PublicationPageInfo, 0, len(result.Records))
	for _, rec := range result.Records {
		var val documentValue
		if err := json.Unmarshal(rec.Value, &val); err != nil || val.Path == "" {
			continue
		}
		info := PublicationPageInfo{
			Title:       val.Title,
			Description: val.Description,
			Link:        pubBase + val.Path,
			ATURI:       rec.URI,
		}
		if t, err := time.Parse(time.RFC3339, val.PublishedAt); err == nil {
			info.Date = t
		}
		docs = append(docs, info)
	}

	sort.Slice(docs, func(i, j int) bool {
		return docs[i].Date.After(docs[j].Date)
	})

	return docs, nil
}
