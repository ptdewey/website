package atproto

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ptdewey/cedar/internal/config"
	"github.com/ptdewey/cedar/internal/parser"
	libleaflet "github.com/ptdewey/standard-site-go/leaflet"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

func convertMarkdown(md string) *libleaflet.Document {
	src := []byte(md)
	p := goldmark.New(goldmark.WithExtensions(extension.GFM, extension.Footnote)).Parser()
	return libleaflet.Convert(src, p.Parse(text.NewReader(src)))
}

// makePage builds a minimal parser.Page for use in tests.
// meta is merged into the default metadata (title, slug).
func makePage(title, slug, rawMD string, meta ...map[string]any) parser.Page {
	m := map[string]any{"title": title, "slug": slug}
	for _, extra := range meta {
		for k, v := range extra {
			m[k] = v
		}
	}
	return parser.Page{
		Metadata:    m,
		RawMarkdown: rawMD,
		Date:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Route: &config.Route{
			ContentPath:   "posts",
			OutputPattern: "/posts/:slug",
		},
	}
}

func testCfg() *config.Config { return &config.Config{ContentDir: "content"} }

func TestBuildDocumentRecordContent(t *testing.T) {
	page := parser.Page{
		Metadata:    map[string]any{"title": "Test Post", "slug": "test-post"},
		RawMarkdown: "# Test\n\nSome **bold** text.",
		PlainText:   "Test Some bold text.",
		Date:        time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC),
		Route: &config.Route{
			ContentPath:   "posts",
			OutputPattern: "/posts/:slug",
		},
	}
	leafletDoc := convertMarkdown(page.RawMarkdown)

	doc := buildDocumentRecord(page, "at://did:plc:abc/site.standard.publication/xyz", testCfg(), leafletDoc, "")

	if doc.Type != "site.standard.document" {
		t.Errorf("$type = %q, want %q", doc.Type, "site.standard.document")
	}
	if doc.Title != "Test Post" {
		t.Errorf("title = %q, want %q", doc.Title, "Test Post")
	}
	if doc.Content == nil {
		t.Fatal("content is nil")
	}
	leafletContent, ok := doc.Content.(*libleaflet.Document)
	if !ok {
		t.Fatalf("content is not *leaflet.Document: %T", doc.Content)
	}
	if leafletContent.Type != "pub.leaflet.content" {
		t.Errorf("content.$type = %q, want %q", leafletContent.Type, "pub.leaflet.content")
	}
	if len(leafletContent.Pages) == 0 {
		t.Error("content.pages is empty")
	}
	if doc.TextContent != page.PlainText {
		t.Errorf("textContent = %q, want %q", doc.TextContent, page.PlainText)
	}
}

func TestBuildDocumentRecordContentNested(t *testing.T) {
	// Verify the full JSON structure has $type at both levels
	page := makePage("Nested", "nested", "hello")
	doc := buildDocumentRecord(page, "at://did:plc:abc/site.standard.publication/xyz", testCfg(), convertMarkdown(page.RawMarkdown), "")
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if m["$type"] != "site.standard.document" {
		t.Errorf("top $type = %v", m["$type"])
	}

	content, ok := m["content"].(map[string]any)
	if !ok {
		t.Fatalf("content is not an object: %T", m["content"])
	}
	if content["$type"] != "pub.leaflet.content" {
		t.Errorf("content.$type = %v, want pub.leaflet.content", content["$type"])
	}
}

func TestBuildDocumentRecordTags(t *testing.T) {
	page := makePage("Tagged", "tagged", "content", map[string]any{"tags": []any{"go", "atproto", "leaflet"}})
	doc := buildDocumentRecord(page, "at://did:plc:abc/site.standard.publication/xyz", testCfg(), nil, "")

	if len(doc.Tags) != 3 {
		t.Fatalf("tags = %d, want 3", len(doc.Tags))
	}
	expected := []string{"go", "atproto", "leaflet"}
	for i, tag := range doc.Tags {
		if tag != expected[i] {
			t.Errorf("tag[%d] = %q, want %q", i, tag, expected[i])
		}
	}
}

func TestBuildDocumentRecordDescription(t *testing.T) {
	page := makePage("Described", "described", "content", map[string]any{"description": "A test post about things."})
	doc := buildDocumentRecord(page, "at://did:plc:abc/site.standard.publication/xyz", testCfg(), nil, "")

	if doc.Description != "A test post about things." {
		t.Errorf("description = %q, want %q", doc.Description, "A test post about things.")
	}
}

func TestBuildDocumentRecordNoDescription(t *testing.T) {
	page := makePage("No Desc", "no-desc", "content")
	doc := buildDocumentRecord(page, "at://did:plc:abc/site.standard.publication/xyz", testCfg(), nil, "")

	if doc.Description != "" {
		t.Errorf("description = %q, want empty", doc.Description)
	}
}

func TestBuildDocumentRecordSite(t *testing.T) {
	page := makePage("Test", "test", "content")
	pubURI := "at://did:plc:abc123/site.standard.publication/rkey456"
	doc := buildDocumentRecord(page, pubURI, testCfg(), nil, "")

	if doc.Site != pubURI {
		t.Errorf("site = %q, want %q", doc.Site, pubURI)
	}
}

func TestBuildDocumentRecordPublishedAt(t *testing.T) {
	date := time.Date(2026, 2, 21, 15, 30, 0, 0, time.UTC)
	page := makePage("Test", "test", "content")
	page.Date = date
	doc := buildDocumentRecord(page, "at://did:plc:abc/site.standard.publication/xyz", testCfg(), nil, "")

	if doc.PublishedAt != "2026-02-21T15:30:00Z" {
		t.Errorf("publishedAt = %q, want %q", doc.PublishedAt, "2026-02-21T15:30:00Z")
	}
}

func TestBuildPublicationRecord(t *testing.T) {
	pubCfg := config.Publication{
		Name:        "My Blog",
		URL:         "https://example.com",
		Description: "A test blog",
	}

	pub := buildPublicationRecord(pubCfg)

	if pub.Type != "site.standard.publication" {
		t.Errorf("$type = %q, want %q", pub.Type, "site.standard.publication")
	}
	if pub.Name != "My Blog" {
		t.Errorf("name = %q, want %q", pub.Name, "My Blog")
	}
	if pub.URL != "https://example.com" {
		t.Errorf("url = %q, want %q", pub.URL, "https://example.com")
	}
	if pub.Description == nil || *pub.Description != "A test blog" {
		t.Errorf("description = %v, want %q", pub.Description, "A test blog")
	}
	if pub.Preferences == nil || !pub.Preferences.ShowInDiscover {
		t.Errorf("preferences.showInDiscover = %v, want true (default)", pub.Preferences)
	}
}

func TestBuildPublicationRecordNoDescription(t *testing.T) {
	pubCfg := config.Publication{
		Name: "My Blog",
		URL:  "https://example.com",
	}

	pub := buildPublicationRecord(pubCfg)

	if pub.Description != nil {
		t.Errorf("description = %v, want nil", pub.Description)
	}
}

func TestBuildPublicationRecordShowInDiscoverExplicitFalse(t *testing.T) {
	show := false
	pubCfg := config.Publication{
		Name:           "My Blog",
		URL:            "https://example.com",
		ShowInDiscover: &show,
	}

	pub := buildPublicationRecord(pubCfg)

	if pub.Preferences == nil || pub.Preferences.ShowInDiscover {
		t.Errorf("preferences.showInDiscover = %v, want false", pub.Preferences)
	}
}
