package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ptdewey/cedar/internal/config"
)

func TestProcessDirectorySingleFile(t *testing.T) {
	dir := t.TempDir()
	contentDir := filepath.Join(dir, "content")
	os.MkdirAll(contentDir, 0755)

	os.WriteFile(filepath.Join(contentDir, "index.md"), []byte(`---
title: Home
---
Welcome home.`), 0644)

	cfg := &config.Config{
		ContentDir: contentDir,
		Routes: []config.Route{
			{ContentPath: "index.md", OutputPattern: "/", Template: "_index.html"},
		},
	}

	pages, err := ProcessDirectory(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) != 1 {
		t.Fatalf("len(pages) = %d, want 1", len(pages))
	}
	if pages[0].Metadata["title"] != "Home" {
		t.Errorf("title = %v, want %q", pages[0].Metadata["title"], "Home")
	}
}

func TestProcessDirectoryMultipleFiles(t *testing.T) {
	dir := t.TempDir()
	contentDir := filepath.Join(dir, "content")
	blogDir := filepath.Join(contentDir, "blog")
	os.MkdirAll(blogDir, 0755)

	os.WriteFile(filepath.Join(blogDir, "post1.md"), []byte(`---
title: Post 1
date: 2024-01-01
---
First post.`), 0644)

	os.WriteFile(filepath.Join(blogDir, "post2.md"), []byte(`---
title: Post 2
date: 2024-02-01
---
Second post.`), 0644)

	cfg := &config.Config{
		ContentDir: contentDir,
		Routes: []config.Route{
			{ContentPath: "blog", OutputPattern: "/blog/:slug", Template: "post.html"},
		},
	}

	pages, err := ProcessDirectory(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) != 2 {
		t.Fatalf("len(pages) = %d, want 2", len(pages))
	}
}

func TestProcessDirectoryMissingPath(t *testing.T) {
	dir := t.TempDir()
	contentDir := filepath.Join(dir, "content")
	os.MkdirAll(contentDir, 0755)

	cfg := &config.Config{
		ContentDir: contentDir,
		Routes: []config.Route{
			{ContentPath: "nonexistent", OutputPattern: "/", Template: "x.html"},
		},
	}

	_, err := ProcessDirectory(cfg)
	if err == nil {
		t.Fatal("expected error for non-existent content path")
	}
}

func TestProcessDirectoryEmptyDir(t *testing.T) {
	dir := t.TempDir()
	contentDir := filepath.Join(dir, "content")
	emptyDir := filepath.Join(contentDir, "empty")
	os.MkdirAll(emptyDir, 0755)

	cfg := &config.Config{
		ContentDir: contentDir,
		Routes: []config.Route{
			{ContentPath: "empty", OutputPattern: "/empty/:slug", Template: "x.html"},
		},
	}

	_, err := ProcessDirectory(cfg)
	if err == nil {
		t.Fatal("expected error for empty directory")
	}
}

func TestProcessDirectoryNonMarkdownFile(t *testing.T) {
	dir := t.TempDir()
	contentDir := filepath.Join(dir, "content")
	os.MkdirAll(contentDir, 0755)

	os.WriteFile(filepath.Join(contentDir, "readme.txt"), []byte("not markdown"), 0644)

	cfg := &config.Config{
		ContentDir: contentDir,
		Routes: []config.Route{
			{ContentPath: "readme.txt", OutputPattern: "/", Template: "x.html"},
		},
	}

	_, err := ProcessDirectory(cfg)
	if err == nil {
		t.Fatal("expected error for non-markdown file")
	}
}

func TestProcessDirectoryMultipleRoutes(t *testing.T) {
	dir := t.TempDir()
	contentDir := filepath.Join(dir, "content")
	os.MkdirAll(filepath.Join(contentDir, "blog"), 0755)

	os.WriteFile(filepath.Join(contentDir, "index.md"), []byte(`---
title: Home
---
Welcome.`), 0644)

	os.WriteFile(filepath.Join(contentDir, "blog", "post.md"), []byte(`---
title: A Post
---
Content.`), 0644)

	cfg := &config.Config{
		ContentDir: contentDir,
		Routes: []config.Route{
			{ContentPath: "index.md", OutputPattern: "/", Template: "_index.html"},
			{ContentPath: "blog", OutputPattern: "/blog/:slug", Template: "post.html"},
		},
	}

	pages, err := ProcessDirectory(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) != 2 {
		t.Fatalf("len(pages) = %d, want 2", len(pages))
	}
}
