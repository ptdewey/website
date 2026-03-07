package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseToml(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cedar.toml")
	os.WriteFile(path, []byte(`
publish_dir = "dist"
static_dir = "assets"
content_dir = "posts"
template_dir = "tmpl"
template_ext = ".html"
clean_build = true
build_draft = true
build_future = true
allow_unsafe_html = true
copyright = "2024 Test"

[rss]
generate = true
title = "My Blog"
description = "A test blog"
url = "example.com"

[atproto]
handle = "test.bsky.social"

[atproto.publications.blog]
name = "Test Pub"
url = "https://example.com"

[[routes]]
content_path = "blog"
output_pattern = "/blog/:slug"
template = "post.html"
generate_rss = true
publish = "blog"
`), 0644)

	cfg, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.PublishDir != "dist" {
		t.Errorf("PublishDir = %q, want %q", cfg.PublishDir, "dist")
	}
	if cfg.StaticDir != "assets" {
		t.Errorf("StaticDir = %q, want %q", cfg.StaticDir, "assets")
	}
	if cfg.ContentDir != "posts" {
		t.Errorf("ContentDir = %q, want %q", cfg.ContentDir, "posts")
	}
	if cfg.TemplateDir != "tmpl" {
		t.Errorf("TemplateDir = %q, want %q", cfg.TemplateDir, "tmpl")
	}
	if cfg.TemplateExt != ".html" {
		t.Errorf("TemplateExt = %q, want %q", cfg.TemplateExt, ".html")
	}
	if !cfg.CleanBuild {
		t.Error("CleanBuild should be true")
	}
	if !cfg.BuildDraft {
		t.Error("BuildDraft should be true")
	}
	if !cfg.BuildFuture {
		t.Error("BuildFuture should be true")
	}
	if !cfg.AllowUnsafeHTML {
		t.Error("AllowUnsafeHTML should be true")
	}
	if cfg.Copyright != "2024 Test" {
		t.Errorf("Copyright = %q, want %q", cfg.Copyright, "2024 Test")
	}

	// RSS
	if !cfg.RSS.Generate {
		t.Error("RSS.Generate should be true")
	}
	if cfg.RSS.Title != "My Blog" {
		t.Errorf("RSS.Title = %q, want %q", cfg.RSS.Title, "My Blog")
	}

	// ATProto
	if cfg.ATProto.Handle != "test.bsky.social" {
		t.Errorf("ATProto.Handle = %q, want %q", cfg.ATProto.Handle, "test.bsky.social")
	}

	// Routes
	if len(cfg.Routes) != 1 {
		t.Fatalf("len(Routes) = %d, want 1", len(cfg.Routes))
	}
	r := cfg.Routes[0]
	if r.ContentPath != "blog" {
		t.Errorf("Route.ContentPath = %q, want %q", r.ContentPath, "blog")
	}
	if r.OutputPattern != "/blog/:slug" {
		t.Errorf("Route.OutputPattern = %q, want %q", r.OutputPattern, "/blog/:slug")
	}
	if !r.GenerateRSS {
		t.Error("Route.GenerateRSS should be true")
	}
	if r.Publish != "blog" {
		t.Errorf("Route.Publish = %q, want %q", r.Publish, "blog")
	}

	// Publications
	if len(cfg.ATProto.Publications) != 1 {
		t.Fatalf("len(Publications) = %d, want 1", len(cfg.ATProto.Publications))
	}
	pub, ok := cfg.ATProto.Publications["blog"]
	if !ok {
		t.Fatal("missing publication 'blog'")
	}
	if pub.Name != "Test Pub" {
		t.Errorf("Publication.Name = %q, want %q", pub.Name, "Test Pub")
	}
}

func TestParseJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cedar.json")
	os.WriteFile(path, []byte(`{
		"publish_dir": "out",
		"content_dir": "src",
		"routes": [
			{"content_path": "pages", "output_pattern": "/:slug", "template": "page.html"}
		]
	}`), 0644)

	cfg, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PublishDir != "out" {
		t.Errorf("PublishDir = %q, want %q", cfg.PublishDir, "out")
	}
	if cfg.ContentDir != "src" {
		t.Errorf("ContentDir = %q, want %q", cfg.ContentDir, "src")
	}
	if len(cfg.Routes) != 1 {
		t.Fatalf("len(Routes) = %d, want 1", len(cfg.Routes))
	}
}

func TestParseYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cedar.yaml")
	os.WriteFile(path, []byte(`
publish_dir: build
content_dir: md
routes:
  - content_path: articles
    output_pattern: /articles/:slug
    template: article.html
`), 0644)

	cfg, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PublishDir != "build" {
		t.Errorf("PublishDir = %q, want %q", cfg.PublishDir, "build")
	}
	if len(cfg.Routes) != 1 {
		t.Fatalf("len(Routes) = %d, want 1", len(cfg.Routes))
	}
}

func TestParseDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cedar.toml")
	os.WriteFile(path, []byte(""), 0644)

	cfg, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PublishDir != "public" {
		t.Errorf("PublishDir = %q, want %q", cfg.PublishDir, "public")
	}
	if cfg.StaticDir != "static" {
		t.Errorf("StaticDir = %q, want %q", cfg.StaticDir, "static")
	}
	if cfg.ContentDir != "content" {
		t.Errorf("ContentDir = %q, want %q", cfg.ContentDir, "content")
	}
	if cfg.TemplateDir != "templates" {
		t.Errorf("TemplateDir = %q, want %q", cfg.TemplateDir, "templates")
	}
	if cfg.TemplateExt != ".tmpl" {
		t.Errorf("TemplateExt = %q, want %q", cfg.TemplateExt, ".tmpl")
	}
	if cfg.CleanBuild {
		t.Error("CleanBuild should default to false")
	}
	if cfg.RSS.Title != "Your Site" {
		t.Errorf("RSS.Title = %q, want %q", cfg.RSS.Title, "Your Site")
	}
}

func TestParseInvalidExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cedar.xml")
	os.WriteFile(path, []byte("<config/>"), 0644)

	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for invalid extension")
	}
}

func TestParseFileNotFound(t *testing.T) {
	_, err := Parse("/nonexistent/cedar.toml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestParseInvalidToml(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cedar.toml")
	os.WriteFile(path, []byte("{{invalid toml"), 0644)

	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for invalid TOML")
	}
}
