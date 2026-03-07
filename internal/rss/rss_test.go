package rss

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ptdewey/cedar/internal/config"
	"github.com/ptdewey/cedar/internal/parser"
)

func TestGenerateRSS(t *testing.T) {
	dir := t.TempDir()

	cfg := &config.Config{
		PublishDir: dir,
		RSS: config.RSS{
			Generate:    true,
			Title:       "Test Blog",
			Description: "A test blog",
			URL:         "example.com",
		},
	}

	route := &config.Route{
		ContentPath:   "blog",
		OutputPattern: "/blog/:slug",
		Template:      "post.html",
		GenerateRSS:   true,
	}

	pages := []parser.Page{
		{
			Metadata: map[string]any{
				"title":       "First Post",
				"slug":        "first-post",
				"description": "The first post",
				"categories":  []any{"tech", "go"},
			},
			Content: "<p>First post content</p>",
			Route:   route,
			Date:    time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			Metadata: map[string]any{
				"title": "Second Post",
				"slug":  "second-post",
			},
			Content: "<p>Second post content</p>",
			Route:   route,
			Date:    time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	if err := GenerateRSS(pages, cfg); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "rss.xml"))
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)
	if !strings.Contains(content, "<?xml") {
		t.Error("missing XML header")
	}

	// Parse and validate structure
	var rss RSS
	if err := xml.Unmarshal(data[len(xml.Header):], &rss); err != nil {
		// try with header
		if err := xml.Unmarshal(data, &rss); err != nil {
			t.Fatalf("invalid RSS XML: %v", err)
		}
	}

	if rss.Version != "2.0" {
		t.Errorf("RSS version = %q, want %q", rss.Version, "2.0")
	}
	if rss.Channel.Title != "Test Blog" {
		t.Errorf("Channel.Title = %q, want %q", rss.Channel.Title, "Test Blog")
	}
	if rss.Channel.Description != "A test blog" {
		t.Errorf("Channel.Description = %q, want %q", rss.Channel.Description, "A test blog")
	}
	if len(rss.Channel.Items) != 2 {
		t.Fatalf("len(Items) = %d, want 2", len(rss.Channel.Items))
	}

	// Check first item
	item := rss.Channel.Items[0]
	if item.Title != "First Post" {
		t.Errorf("Item.Title = %q, want %q", item.Title, "First Post")
	}
	if !strings.Contains(item.Link, "first-post") {
		t.Errorf("Item.Link = %q, should contain %q", item.Link, "first-post")
	}
	if item.Description != "The first post" {
		t.Errorf("Item.Description = %q, want %q", item.Description, "The first post")
	}
	if item.Category != "tech, go" {
		t.Errorf("Item.Category = %q, want %q", item.Category, "tech, go")
	}
}

func TestGenerateRSSSkipsNonRSSRoutes(t *testing.T) {
	dir := t.TempDir()

	cfg := &config.Config{
		PublishDir: dir,
		RSS: config.RSS{
			Generate:    true,
			Title:       "Test",
			Description: "Test",
			URL:         "example.com",
		},
	}

	rssRoute := &config.Route{
		ContentPath: "blog",
		OutputPattern: "/blog/:slug",
		GenerateRSS: true,
	}
	noRSSRoute := &config.Route{
		ContentPath: "pages",
		OutputPattern: "/:slug",
		GenerateRSS: false,
	}

	pages := []parser.Page{
		{
			Metadata: map[string]any{"title": "Blog Post", "slug": "post"},
			Content:  "<p>blog</p>",
			Route:    rssRoute,
			Date:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Metadata: map[string]any{"title": "About", "slug": "about"},
			Content:  "<p>about</p>",
			Route:    noRSSRoute,
			Date:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Metadata: map[string]any{"title": "Orphan", "slug": "orphan"},
			Content:  "<p>no route</p>",
			Route:    nil,
			Date:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	if err := GenerateRSS(pages, cfg); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(filepath.Join(dir, "rss.xml"))
	var rss RSS
	xml.Unmarshal(data, &rss)

	if len(rss.Channel.Items) != 1 {
		t.Fatalf("len(Items) = %d, want 1 (only RSS-enabled route)", len(rss.Channel.Items))
	}
	if rss.Channel.Items[0].Title != "Blog Post" {
		t.Errorf("Item.Title = %q, want %q", rss.Channel.Items[0].Title, "Blog Post")
	}
}

func TestGenerateRSSNoPages(t *testing.T) {
	dir := t.TempDir()

	cfg := &config.Config{
		PublishDir: dir,
		RSS: config.RSS{
			Generate:    true,
			Title:       "Empty",
			Description: "No posts",
			URL:         "example.com",
		},
	}

	if err := GenerateRSS(nil, cfg); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "rss.xml"))
	if err != nil {
		t.Fatal(err)
	}

	var rss RSS
	xml.Unmarshal(data, &rss)
	if len(rss.Channel.Items) != 0 {
		t.Errorf("len(Items) = %d, want 0", len(rss.Channel.Items))
	}
}

func TestGenerateRSSLinkFormat(t *testing.T) {
	dir := t.TempDir()

	cfg := &config.Config{
		PublishDir: dir,
		RSS: config.RSS{
			Generate: true,
			Title:    "Test",
			URL:      "myblog.com/",
		},
	}

	route := &config.Route{
		ContentPath:   "posts",
		OutputPattern: "/posts/:slug",
		GenerateRSS:   true,
	}

	pages := []parser.Page{
		{
			Metadata: map[string]any{"title": "Test", "slug": "my-post"},
			Content:  "<p>x</p>",
			Route:    route,
			Date:     time.Now(),
		},
	}

	GenerateRSS(pages, cfg)

	data, _ := os.ReadFile(filepath.Join(dir, "rss.xml"))
	var rss RSS
	xml.Unmarshal(data, &rss)

	link := rss.Channel.Items[0].Link
	if !strings.HasPrefix(link, "https://") {
		t.Errorf("link should start with https://, got %q", link)
	}
	if strings.Contains(link, "index.html") {
		t.Errorf("link should not contain index.html, got %q", link)
	}
}
