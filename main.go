package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	chromahtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v2"
)

type Page struct {
	Metadata map[string]interface{} `json:"metadata"`
	Content  string                 `json:"content"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description,omitempty"`
	PubDate     string `xml:"pubDate"`
	Category    string `xml:"category,omitempty"`
	Content     string `xml:"content"`
}

func generateSlug(title string) string {
	slug := strings.ToLower(strings.TrimSpace(title))
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, ".", "")
	return slug
}

func parseFrontMatter(content []byte) (map[string]interface{}, []byte, error) {
	contentStr := string(content)
	if !strings.HasPrefix(contentStr, "---") {
		return nil, content, nil
	}
	parts := strings.SplitN(contentStr, "---", 3)
	if len(parts) < 3 {
		return nil, content, fmt.Errorf("invalid front-matter format")
	}
	var metadata map[string]interface{}
	if err := yaml.Unmarshal([]byte(parts[1]), &metadata); err != nil {
		return nil, nil, err
	}
	return metadata, []byte(parts[2]), nil
}

func getReadingTime(text string) int {
	words := strings.Fields(text)
	wordCount := len(words)

	// reading/speaking rate
	wordsPerMinute := 200.0
	return int(math.Round(float64(wordCount) / wordsPerMinute))
}

func processMarkdownFile(path string) (Page, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Page{}, err
	}

	metadata, markdownContent, err := parseFrontMatter(content)
	if err != nil {
		return Page{}, err
	}

	if _, ok := metadata["slug"]; !ok {
		fname := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		metadata["slug"] = generateSlug(fname)
	}

	var htmlContent bytes.Buffer

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			highlighting.NewHighlighting(
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
					chromahtml.WithClasses(true),
				),
				// highlighting.WithStyle("monokai"),
				// highlighting.WithCustomStyle(darkearth),
			),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	if err := md.Convert(markdownContent, &htmlContent); err != nil {
		return Page{}, err
	}

	metadata["read_time"] = getReadingTime(string(markdownContent))

	return Page{
		Metadata: metadata,
		Content:  htmlContent.String(),
	}, nil
}

func processDirectory(dir string) ([]Page, error) {
	var writings []Page

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), ".md") {
			page, err := processMarkdownFile(path)
			if err != nil {
				return err
			}

			writings = append(writings, page)
		}
		return nil
	})

	return writings, err
}

func generateRSSFromJSON(inputPath string, outputPath string) error {
	file, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	var pages []Page
	if err := json.Unmarshal(file, &pages); err != nil {
		return err
	}

	channel := Channel{
		Title:       "Patrick Dewey's Personal Site",
		Link:        "https://pdewey.com/writing",
		Description: "RSS feed of Patrick Dewey's website",
		PubDate:     time.Now().Format(time.RFC1123Z),
	}

	for _, page := range pages {
		var categories []string
		if rawCategories, ok := page.Metadata["categories"].([]interface{}); ok {
			for _, category := range rawCategories {
				if strCategory, ok := category.(string); ok {
					categories = append(categories, strCategory)
				}
			}
		}

		description, ok := page.Metadata["description"].(string)
		if !ok {
			description = ""
		}

		item := Item{
			Title:       page.Metadata["title"].(string),
			Link:        fmt.Sprintf("https://pdewey.com/writing/%s", page.Metadata["slug"].(string)),
			Description: description,
			Content:     page.Content,
			PubDate:     page.Metadata["date"].(string),
			Category:    strings.Join(categories, ", "),
		}
		channel.Items = append(channel.Items, item)
	}

	rss := RSS{
		Version: "2.0",
		Channel: channel,
	}

	output, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return err
	}

	rssHeader := []byte(xml.Header)
	output = append(rssHeader, output...)

	return os.WriteFile(outputPath, output, 0644)
}

func writeJSONFile(data interface{}, outputPath string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, jsonData, 0644)
}

func main() {
	writings, err := processDirectory("content")
	if err != nil {
		fmt.Printf("Error processing directory: %v\n", err)
		os.Exit(1)
	}

	sort.Slice(writings, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", writings[i].Metadata["date"].(string))
		dateJ, _ := time.Parse("2006-01-02", writings[j].Metadata["date"].(string))
		return dateI.After(dateJ)
	})

	// TODO: spread out output to different files based on subdir structure within content
	if err := writeJSONFile(writings, "static/data/writing.json"); err != nil {
		fmt.Printf("Error writing writing.json: %v\n", err)
		os.Exit(1)
	}

	if err := generateRSSFromJSON("static/data/writing.json", "static/rss.xml"); err != nil {
		fmt.Printf("Error writing rss.xml: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully generated writing.json and rss.xml")
}
