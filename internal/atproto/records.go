package atproto

import (
	"time"

	"github.com/ptdewey/cedar/internal/config"
	"github.com/ptdewey/cedar/internal/parser"
	"github.com/ptdewey/standard-site-go/standard"
)

// markdownContent is the content value for site.standard.content.markdown documents.
type markdownContent struct {
	Type string `json:"$type"`
	Text string `json:"text"`
}

func boolPref(p *bool) bool {
	return p == nil || *p
}

func convertColor(c *standard.ThemeColor) *standard.ThemeColor {
	if c == nil {
		return nil
	}
	return &standard.ThemeColor{Type: "site.standard.theme.color#rgb", R: c.R, G: c.G, B: c.B}
}

func buildPublicationRecord(pub config.Publication) *standard.Publication {
	rec := &standard.Publication{
		Type: "site.standard.publication",
		URL:  pub.URL,
		Name: pub.Name,
		Preferences: &standard.PublicationPreferences{
			ShowInDiscover: boolPref(pub.ShowInDiscover),
		},
	}
	if pub.Description != "" {
		rec.Description = &pub.Description
	}
	if t := pub.BasicTheme; t != nil {
		rec.BasicTheme = &standard.BasicTheme{
			Type:             "site.standard.theme.basic",
			Accent:           convertColor(t.Accent),
			Background:       convertColor(t.Background),
			Foreground:       convertColor(t.Foreground),
			AccentForeground: convertColor(t.AccentForeground),
		}
	}
	return rec
}

// documentPath returns the path value for a document record based on the
// publication's configured path mode. For "slug" mode, the page's slug is
// used. For "rkey" mode (or any other value), the rkey is used. When rkey
// is empty (new record, rkey not yet known), an empty string is returned
// so the caller can set it after creation.
func documentPath(mode string, page parser.Page, rkey string) string {
	if mode == "slug" {
		if slug, ok := page.Metadata["slug"].(string); ok {
			return "/" + slug
		}
	}
	if rkey != "" {
		return "/" + rkey
	}
	return ""
}

func buildDocumentRecord(page parser.Page, publicationURI string, cfg *config.Config, content any, path string) *standard.Document {
	title, _ := page.Metadata["title"].(string)

	opts := []standard.DocumentOption{
		standard.WithUpdatedAt(time.Now()),
		standard.WithContent(content),
	}
	if path != "" {
		opts = append(opts, standard.WithPath(path))
	}
	if desc, ok := page.Metadata["description"].(string); ok {
		opts = append(opts, standard.WithDescription(desc))
	}
	if raw, ok := page.Metadata["tags"].([]any); ok {
		tags := make([]string, 0, len(raw))
		for _, t := range raw {
			if s, ok := t.(string); ok {
				tags = append(tags, s)
			}
		}
		opts = append(opts, standard.WithTags(tags...))
	}

	return standard.NewDocument(title, publicationURI, page.Date, opts...)
}
