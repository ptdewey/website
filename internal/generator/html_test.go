package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ptdewey/cedar/internal/config"
	"github.com/ptdewey/cedar/internal/parser"
)

// TestWriteHTMLFiles_noPublicationFetch verifies that when no publications have
// include_in_build=true, the build succeeds without any network calls and
// PublicationPages is available (empty) in templates.
func TestWriteHTMLFiles_noPublicationFetch(t *testing.T) {
	dir := t.TempDir()
	templateDir := filepath.Join(dir, "templates")
	partialsDir := filepath.Join(templateDir, "partials")
	if err := os.MkdirAll(partialsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// ParseGlob requires at least one match; provide an empty partial.
	if err := os.WriteFile(filepath.Join(partialsDir, "noop.tmpl"), []byte(`{{define "noop"}}{{end}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	// Template that iterates PublicationPages — must not panic on empty map.
	tmplContent := `{{define "page.tmpl"}}{{range $k, $v := .PublicationPages}}{{$k}}:{{len $v}} {{end}}{{end}}`
	if err := os.WriteFile(filepath.Join(templateDir, "page.tmpl"), []byte(tmplContent), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		PublishDir:  filepath.Join(dir, "public"),
		TemplateDir: templateDir,
		TemplateExt: ".tmpl",
		ContentDir:  "content",
		ATProto: config.ATProto{
			Handle: "", // no handle → no DID/PDS resolution → no fetch
			Publications: map[string]config.Publication{
				"blog": {
					Name:           "Blog",
					URL:            "https://example.pub",
					IncludeInBuild: false,
				},
			},
		},
	}

	route := &config.Route{
		ContentPath:   "posts",
		OutputPattern: "/posts/:slug",
		Template:      "page.tmpl",
	}
	pages := []parser.Page{
		{
			Metadata:   map[string]any{"title": "Test", "slug": "test"},
			Content:    "<p>hi</p>",
			SourcePath: "content/posts/test.md",
			Route:      route,
		},
	}

	if err := WriteHTMLFiles(pages, cfg.PublishDir, cfg); err != nil {
		t.Fatalf("WriteHTMLFiles returned error: %v", err)
	}

	out := filepath.Join(cfg.PublishDir, "posts", "test", "index.html")
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("expected output file %s: %v", out, err)
	}
}

// TestWriteHTMLFiles_publicationPagesSkippedWhenIncludeInBuildFalse verifies
// that a publication with include_in_build=false is not included in PublicationPages.
func TestWriteHTMLFiles_publicationPagesSkippedWhenIncludeInBuildFalse(t *testing.T) {
	dir := t.TempDir()
	templateDir := filepath.Join(dir, "templates")
	partialsDir := filepath.Join(templateDir, "partials")
	if err := os.MkdirAll(partialsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(partialsDir, "noop.tmpl"), []byte(`{{define "noop"}}{{end}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	// Write the count of publication keys into the output so we can assert.
	tmplContent := `{{define "page.tmpl"}}keys:{{len .PublicationPages}}{{end}}`
	if err := os.WriteFile(filepath.Join(templateDir, "page.tmpl"), []byte(tmplContent), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		PublishDir:  filepath.Join(dir, "public"),
		TemplateDir: templateDir,
		TemplateExt: ".tmpl",
		ContentDir:  "content",
		ATProto: config.ATProto{
			Handle: "",
			Publications: map[string]config.Publication{
				"blog": {Name: "Blog", URL: "https://example.pub", IncludeInBuild: false},
			},
		},
	}

	pages := []parser.Page{
		{
			Metadata:   map[string]any{"title": "T", "slug": "t"},
			Content:    "",
			SourcePath: "content/posts/t.md",
			Route: &config.Route{
				ContentPath:   "posts",
				OutputPattern: "/posts/:slug",
				Template:      "page.tmpl",
			},
		},
	}

	if err := WriteHTMLFiles(pages, cfg.PublishDir, cfg); err != nil {
		t.Fatalf("WriteHTMLFiles error: %v", err)
	}

	out := filepath.Join(cfg.PublishDir, "posts", "t", "index.html")
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}
	if string(data) != "keys:0" {
		t.Errorf("output = %q, want %q", string(data), "keys:0")
	}
}
