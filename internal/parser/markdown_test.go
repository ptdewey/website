package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ptdewey/cedar/internal/config"
)

func TestParseFrontMatter(t *testing.T) {
	input := []byte(`---
title: Hello World
date: 2024-01-15
tags:
  - go
  - test
---
# Content here`)

	meta, body, err := parseFrontMatter(input)
	if err != nil {
		t.Fatal(err)
	}
	if meta["title"] != "Hello World" {
		t.Errorf("title = %v, want %q", meta["title"], "Hello World")
	}
	if meta["date"] != "2024-01-15" {
		t.Errorf("date = %v, want %q", meta["date"], "2024-01-15")
	}
	if !strings.Contains(string(body), "# Content here") {
		t.Errorf("body should contain markdown content, got %q", string(body))
	}
}

func TestParseFrontMatterNoFrontMatter(t *testing.T) {
	input := []byte("# Just markdown\nNo front matter here.")
	meta, body, err := parseFrontMatter(input)
	if err != nil {
		t.Fatal(err)
	}
	if meta != nil {
		t.Errorf("metadata should be nil, got %v", meta)
	}
	if string(body) != string(input) {
		t.Errorf("body should equal input")
	}
}

func TestParseFrontMatterInvalid(t *testing.T) {
	input := []byte("---\nonly one delimiter")
	_, _, err := parseFrontMatter(input)
	if err == nil {
		t.Fatal("expected error for invalid front matter")
	}
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Hello World", "hello-world"},
		{"my.post.title", "myposttitle"},
		{"  Spaces  ", "spaces"},
		{"already-slug", "already-slug"},
		{"UPPER CASE", "upper-case"},
	}

	for _, tt := range tests {
		got := generateSlug(tt.input)
		if got != tt.want {
			t.Errorf("generateSlug(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestGetReadingTime(t *testing.T) {
	// 200 words = 1 minute
	words := strings.Repeat("word ", 200)
	if got := getReadingTime(words); got != 1 {
		t.Errorf("getReadingTime(200 words) = %d, want 1", got)
	}

	// 400 words = 2 minutes
	words = strings.Repeat("word ", 400)
	if got := getReadingTime(words); got != 2 {
		t.Errorf("getReadingTime(400 words) = %d, want 2", got)
	}

	// 0 words = 0 minutes
	if got := getReadingTime(""); got != 0 {
		t.Errorf("getReadingTime(\"\") = %d, want 0", got)
	}
}

func TestStripMarkdown(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "headings",
			input: "## Heading\nsome text",
			want:  "Heading some text",
		},
		{
			name:  "bold and italic",
			input: "**bold** and __also bold__",
			want:  "bold and also bold",
		},
		{
			name:  "links",
			input: "a [link](https://example.com) here",
			want:  "a link here",
		},
		{
			name:  "images",
			input: "text ![alt](img.png) more",
			want:  "text more",
		},
		{
			name:  "code fences",
			input: "before\n```go\nfmt.Println(\"hi\")\n```\nafter",
			want:  "before after",
		},
		{
			name:  "inline code",
			input: "use `fmt.Println` here",
			want:  "use fmt.Println here",
		},
		{
			name:  "blockquote",
			input: "> quoted text",
			want:  "quoted text",
		},
		{
			name:  "strikethrough",
			input: "~~deleted~~ text",
			want:  "deleted text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripMarkdown(tt.input)
			if got != tt.want {
				t.Errorf("stripMarkdown(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMatchRoute(t *testing.T) {
	routes := []config.Route{
		{ContentPath: "index.md", OutputPattern: "/", Template: "_index.html"},
		{ContentPath: "blog", OutputPattern: "/blog/:slug", Template: "post.html"},
		{ContentPath: "pages", OutputPattern: "/:slug", Template: "page.html"},
	}

	tests := []struct {
		filePath string
		wantNil  bool
		wantPath string
	}{
		{"index.md", false, "index.md"},
		{"blog/first-post.md", false, "blog"},
		{"blog/nested/deep.md", false, "blog"},
		{"pages/about.md", false, "pages"},
		{"unknown/file.md", true, ""},
	}

	for _, tt := range tests {
		r := matchRoute(tt.filePath, routes)
		if tt.wantNil {
			if r != nil {
				t.Errorf("matchRoute(%q) = %v, want nil", tt.filePath, r)
			}
		} else {
			if r == nil {
				t.Fatalf("matchRoute(%q) = nil, want route with ContentPath %q", tt.filePath, tt.wantPath)
			}
			if r.ContentPath != tt.wantPath {
				t.Errorf("matchRoute(%q).ContentPath = %q, want %q", tt.filePath, r.ContentPath, tt.wantPath)
			}
		}
	}
}

func TestGetOutputPath(t *testing.T) {
	tests := []struct {
		name      string
		page      Page
		publishDir string
		want      string
	}{
		{
			name: "index route",
			page: Page{
				Metadata: map[string]any{"slug": "index"},
				Route:    &config.Route{OutputPattern: "/"},
			},
			publishDir: "public",
			want:      filepath.Join("public", "index.html"),
		},
		{
			name: "slug route",
			page: Page{
				Metadata: map[string]any{"slug": "hello-world"},
				Route:    &config.Route{OutputPattern: "/blog/:slug"},
			},
			publishDir: "public",
			want:      filepath.Join("public", "blog", "hello-world", "index.html"),
		},
		{
			name: "no route",
			page: Page{
				Metadata: map[string]any{"slug": "orphan"},
				Route:    nil,
			},
			publishDir: "public",
			want:      filepath.Join("public", "orphan", "index.html"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetOutputPath(tt.page, tt.publishDir)
			if got != tt.want {
				t.Errorf("GetOutputPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]any
		wantZero bool
		wantYear int
	}{
		{
			name:     "string date",
			metadata: map[string]any{"date": "2024-01-15"},
			wantYear: 2024,
		},
		{
			name:     "RFC3339",
			metadata: map[string]any{"date": "2024-06-15T10:30:00Z"},
			wantYear: 2024,
		},
		{
			name:     "datetime no tz",
			metadata: map[string]any{"date": "2024-06-15T10:30:00"},
			wantYear: 2024,
		},
		{
			name:     "datetime with space",
			metadata: map[string]any{"date": "2024-06-15 10:30:00"},
			wantYear: 2024,
		},
		{
			name:     "time.Time value",
			metadata: map[string]any{"date": time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)},
			wantYear: 2023,
		},
		{
			name:     "no date",
			metadata: map[string]any{},
			wantZero: true,
		},
		{
			name:     "invalid date string",
			metadata: map[string]any{"date": "not-a-date"},
			wantZero: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseDate(tt.metadata)
			if tt.wantZero {
				if !got.IsZero() {
					t.Errorf("parseDate() = %v, want zero", got)
				}
			} else {
				if got.Year() != tt.wantYear {
					t.Errorf("parseDate().Year() = %d, want %d", got.Year(), tt.wantYear)
				}
			}
		})
	}
}

func TestProcessMarkdownFile(t *testing.T) {
	dir := t.TempDir()
	contentDir := filepath.Join(dir, "content")
	os.MkdirAll(contentDir, 0755)

	mdPath := filepath.Join(contentDir, "test-post.md")
	os.WriteFile(mdPath, []byte(`---
title: Test Post
date: 2024-03-15
description: A test post
tags:
  - testing
---
# Hello

This is a **test** post with a [link](https://example.com).
`), 0644)

	cfg := &config.Config{
		ContentDir: contentDir,
		Routes: []config.Route{
			{ContentPath: "test-post.md", OutputPattern: "/", Template: "_index.html"},
		},
	}

	page, err := ProcessMarkdownFile(mdPath, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if page.Metadata["title"] != "Test Post" {
		t.Errorf("title = %v, want %q", page.Metadata["title"], "Test Post")
	}
	if page.Metadata["slug"] != "test-post" {
		t.Errorf("slug = %v, want %q", page.Metadata["slug"], "test-post")
	}
	if !strings.Contains(page.Content, "<h1>Hello</h1>") {
		t.Errorf("Content should contain <h1>, got %q", page.Content)
	}
	if !strings.Contains(page.Content, "<strong>test</strong>") {
		t.Errorf("Content should contain <strong>, got %q", page.Content)
	}
	if page.RawMarkdown == "" {
		t.Error("RawMarkdown should not be empty")
	}
	if page.PlainText == "" {
		t.Error("PlainText should not be empty")
	}
	if page.Route == nil {
		t.Error("Route should not be nil")
	}
	if page.Date.Year() != 2024 {
		t.Errorf("Date.Year() = %d, want 2024", page.Date.Year())
	}
}

func TestProcessMarkdownFileNoFrontMatter(t *testing.T) {
	dir := t.TempDir()
	contentDir := filepath.Join(dir, "content")
	os.MkdirAll(contentDir, 0755)

	mdPath := filepath.Join(contentDir, "bare.md")
	os.WriteFile(mdPath, []byte("# Just Content\nNo front matter."), 0644)

	cfg := &config.Config{
		ContentDir: contentDir,
		Routes:     []config.Route{},
	}

	page, err := ProcessMarkdownFile(mdPath, cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Should get default metadata
	if page.Metadata["title"] != "bare" {
		t.Errorf("title = %v, want %q", page.Metadata["title"], "bare")
	}
	if page.Metadata["slug"] != "bare" {
		t.Errorf("slug = %v, want %q", page.Metadata["slug"], "bare")
	}
}

func TestProcessMarkdownFileWithCustomSlug(t *testing.T) {
	dir := t.TempDir()
	contentDir := filepath.Join(dir, "content")
	os.MkdirAll(contentDir, 0755)

	mdPath := filepath.Join(contentDir, "post.md")
	os.WriteFile(mdPath, []byte(`---
title: My Post
slug: custom-slug
---
Content`), 0644)

	cfg := &config.Config{
		ContentDir: contentDir,
		Routes:     []config.Route{},
	}

	page, err := ProcessMarkdownFile(mdPath, cfg)
	if err != nil {
		t.Fatal(err)
	}

	if page.Metadata["slug"] != "custom-slug" {
		t.Errorf("slug = %v, want %q", page.Metadata["slug"], "custom-slug")
	}
}
