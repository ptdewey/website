package rss

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ptdewey/cedar/internal/config"
	"github.com/ptdewey/cedar/internal/parser"
)

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

func GenerateRSS(pages []parser.Page, cfg *config.Config) error {
	outputPath := filepath.Join(cfg.PublishDir, "rss.xml")

	channel := Channel{
		Title:       cfg.RSS.Title,
		Link:        cfg.RSS.URL,
		Description: cfg.RSS.Description,
		PubDate:     time.Now().Format(time.RFC1123Z),
	}

	sortedPages := make([]parser.Page, 0, len(pages))
	for _, page := range pages {
		if page.Route != nil && page.Route.GenerateRSS {
			sortedPages = append(sortedPages, page)
		}
	}
	sort.Slice(sortedPages, func(i, j int) bool {
		return sortedPages[i].Date.After(sortedPages[j].Date)
	})

	for _, page := range pages {
		// Skip pages whose routes don't have RSS generation enabled
		if page.Route == nil || !page.Route.GenerateRSS {
			continue
		}

		var categories []string
		if rawCategories, ok := page.Metadata["categories"].([]any); ok {
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

		baseURL := strings.TrimSuffix(cfg.RSS.URL, "/")
		itemPath := strings.TrimPrefix(parser.GetOutputPath(page, ""), "/")
		itemPath = strings.TrimSuffix(itemPath, "/index.html")

		item := Item{
			Title:       page.Metadata["title"].(string),
			Link:        fmt.Sprintf("https://%s/%s", baseURL, itemPath),
			Description: description,
			Content:     page.Content,
			PubDate:     page.Date.Format(time.RFC1123Z),
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
