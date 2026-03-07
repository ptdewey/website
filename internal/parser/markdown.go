package parser

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	chromahtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/goccy/go-yaml"
	"github.com/ptdewey/cedar/internal/config"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
)

type Page struct {
	Metadata    map[string]any `json:"metadata"`
	Content     string         `json:"content"`
	RawMarkdown string         `json:"-"` // Raw markdown (after front matter) for ATProto publishing
	PlainText   string         `json:"-"` // Stripped plain text for ATProto search indexing
	Route       *config.Route  `json:"-"`
	Date        time.Time      `json:"-"` // Parsed date for sorting
	SourcePath  string         `json:"-"` // Path to the source markdown file
}

func ProcessMarkdownFile(path string, cfg *config.Config) (Page, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Page{}, err
	}

	metadata, markdownContent, err := parseFrontMatter(content)
	if err != nil {
		return Page{}, err
	} else if metadata == nil {
		metadata = setDefaultMetadata(path)
	}

	// Match the file to a route
	relPath, err := filepath.Rel(cfg.ContentDir, path)
	if err != nil {
		return Page{}, err
	}

	route := matchRoute(relPath, cfg.Routes)

	// Generate slug if not provided
	if _, ok := metadata["slug"]; !ok {
		fname := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		metadata["slug"] = generateSlug(fname)
	}

	var htmlContent bytes.Buffer

	md := newGoldmarkParser(cfg)
	if err := md.Convert(markdownContent, &htmlContent); err != nil {
		return Page{}, err
	}

	metadata["read_time"] = getReadingTime(string(markdownContent))

	rawMd := strings.TrimSpace(string(markdownContent))

	return Page{
		Metadata:    metadata,
		Content:     htmlContent.String(),
		RawMarkdown: rawMd,
		PlainText:   stripMarkdown(rawMd),
		Route:       route,
		Date:        parseDate(metadata),
		SourcePath:  path,
	}, nil
}

func newGoldmarkParser(cfg *config.Config) goldmark.Markdown {
	opts := []renderer.Option{
		// html.WithHardWraps(),
	}
	if cfg.AllowUnsafeHTML {
		opts = append(opts, html.WithUnsafe())
	}

	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,

			highlighting.NewHighlighting(
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(false),
					chromahtml.WithClasses(true),
				),
			),
		),
		goldmark.WithRendererOptions(opts...),
	)
}

func generateSlug(title string) string {
	slug := strings.ToLower(strings.TrimSpace(title))
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, ".", "")
	return slug
}

func parseFrontMatter(content []byte) (map[string]any, []byte, error) {
	contentStr := string(content)
	// TODO: allow toml metadata
	if !strings.HasPrefix(contentStr, "---") {
		return nil, content, nil
	}
	parts := strings.SplitN(contentStr, "---", 3)
	if len(parts) < 3 {
		return nil, content, fmt.Errorf("invalid front-matter format")
	}
	var metadata map[string]any
	if err := yaml.Unmarshal([]byte(parts[1]), &metadata); err != nil {
		return nil, nil, err
	}
	return metadata, []byte(parts[2]), nil
}

func setDefaultMetadata(path string) map[string]any {
	_, filename := filepath.Split(path)
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	// TODO: allow configuration of default metadata
	return map[string]any{
		"title":       filename,
		"description": "",
		"date":        time.Now().Format("2006-01-02"),
	}
}

func getReadingTime(text string) int {
	words := strings.Fields(text)
	wordCount := len(words)

	// reading/speaking rate
	wordsPerMinute := 200.0
	return int(math.Round(float64(wordCount) / wordsPerMinute))
}

var (
	reImage     = regexp.MustCompile(`!\[[^\]]*\]\([^)]*\)`)
	reLink      = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	reCodeFence = regexp.MustCompile("(?m)^```[\\s\\S]*?^```")
	reHeading   = regexp.MustCompile(`(?m)^#{1,6}\s+`)
)

func stripMarkdown(md string) string {
	s := reCodeFence.ReplaceAllString(md, "")
	s = reImage.ReplaceAllString(s, "")
	s = reLink.ReplaceAllString(s, "$1")
	s = reHeading.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "**", "")
	s = strings.ReplaceAll(s, "__", "")
	s = strings.ReplaceAll(s, "~~", "")
	s = strings.ReplaceAll(s, "`", "")
	s = strings.ReplaceAll(s, "> ", "")
	// Collapse whitespace
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

func matchRoute(filePath string, routes []config.Route) *config.Route {
	filePath = filepath.ToSlash(filePath)

	for _, route := range routes {
		routePath := filepath.ToSlash(route.ContentPath)

		if filePath == routePath {
			return &route
		}

		if strings.HasPrefix(filePath, routePath+"/") {
			return &route
		}
	}

	return nil
}

func GetOutputPath(page Page, publishDir string) string {
	if page.Route == nil {
		slug := page.Metadata["slug"].(string)
		return filepath.Join(publishDir, slug, "index.html")
	}

	outputPath := page.Route.OutputPattern

	if slug, ok := page.Metadata["slug"].(string); ok {
		outputPath = strings.ReplaceAll(outputPath, ":slug", slug)
	}

	if outputPath == "/" {
		outputPath = ""
	} else {
		outputPath = strings.TrimPrefix(outputPath, "/")
	}

	var fileName string
	if outputPath == "" {
		fileName = "index.html"
	} else {
		fileName = filepath.Join(outputPath, "index.html")
	}

	return filepath.Join(publishDir, fileName)
}

func parseDate(metadata map[string]any) time.Time {
	dateVal, ok := metadata["date"]
	if !ok {
		return time.Time{}
	}

	switch v := dateVal.(type) {
	case time.Time:
		return v
	case string:
		formats := []string{
			time.RFC3339,
			"2006-01-02",
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, v); err == nil {
				return t
			}
		}
	}
	return time.Time{}
}
