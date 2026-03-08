package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/pelletier/go-toml/v2"
	"github.com/ptdewey/standard-site-go/standard"
)

// Publication represents a named ATProto publication with its own
// site.standard.publication record.
type Publication struct {
	Name           string               `toml:"name" json:"name" yaml:"name"`
	URL            string               `toml:"url" json:"url" yaml:"url"`
	Description    string               `toml:"description" json:"description" yaml:"description"`
	ShowInDiscover *bool                `toml:"show_in_discover" json:"show_in_discover" yaml:"show_in_discover"`
	ShowComments   *bool                `toml:"show_comments" json:"show_comments" yaml:"show_comments"`
	ShowMentions   *bool                `toml:"show_mentions" json:"show_mentions" yaml:"show_mentions"`
	ShowPrevNext   *bool                `toml:"show_prev_next" json:"show_prev_next" yaml:"show_prev_next"`
	BasicTheme     *standard.BasicTheme `toml:"basic_theme" json:"basic_theme" yaml:"basic_theme"`
	// ContentType controls the document content format used when publishing.
	// Valid values are "markdown" (default) and "leaflet".
	ContentType string `toml:"content_type" json:"content_type" yaml:"content_type"`
	// PathMode controls how the document path is set when publishing.
	// "slug" uses the page slug (e.g. /my-post), "rkey" uses the ATProto
	// record key (default).
	PathMode string `toml:"path_mode" json:"path_mode" yaml:"path_mode"`
	// RKey pins this publication to an existing ATProto record key. When set,
	// cedar will not touch the record and will reuse the existing AT-URI.
	RKey string `toml:"rkey" json:"rkey" yaml:"rkey"`
	// IncludeInBuild controls whether documents from this publication are
	// fetched from the PDS at build time and made available in templates
	// as .PublicationPages.
	IncludeInBuild bool `toml:"include_in_build" json:"include_in_build" yaml:"include_in_build"`
}

type ATProto struct {
	Handle       string                 `toml:"handle" json:"handle" yaml:"handle"`
	Publications map[string]Publication `toml:"publications" json:"publications" yaml:"publications"`
}

type Config struct {
	PublishDir       string  `toml:"publish_dir" json:"publish_dir" yaml:"publish_dir"`
	StaticDir        string  `toml:"static_dir" json:"static_dir" yaml:"static_dir"`
	ContentDir       string  `toml:"content_dir" json:"content_dir" yaml:"content_dir"`
	TemplateDir      string  `toml:"template_dir" json:"template_dir" yaml:"template_dir"`
	TemplateExt      string  `toml:"template_ext" json:"template_ext" yaml:"template_ext"`
	CacheDir         string  `toml:"cache_dir" json:"cache_dir" yaml:"cache_dir"`
	BaseTemplatePath string  `toml:"base_template_path" json:"base_template_path" yaml:"base_template_path"`
	Copyright        string  `toml:"copyright" json:"copyright" yaml:"copyright"`
	CleanBuild       bool    `toml:"clean_build" json:"clean_build" yaml:"clean_build"`
	BuildDraft       bool    `toml:"build_draft" json:"build_draft" yaml:"build_draft"`
	BuildFuture      bool    `toml:"build_future" json:"build_future" yaml:"build_future"`
	AllowUnsafeHTML  bool    `toml:"allow_unsafe_html" json:"allow_unsafe_html" yaml:"allow_unsafe_html"`
	RSS              RSS     `toml:"rss" json:"rss" yaml:"rss"`
	ATProto          ATProto `toml:"atproto" json:"atproto" yaml:"atproto"`
	// TODO: atom feed support
	Routes []Route `toml:"routes" json:"routes" yaml:"routes"`
}

type RSS struct {
	Generate    bool   `toml:"generate" json:"generate" yaml:"generate"`
	Title       string `toml:"title" json:"title" yaml:"title"`
	Description string `toml:"description" json:"description" yaml:"description"`
	URL         string `toml:"url" json:"url" yaml:"url"`
}

type Route struct {
	ContentPath   string `toml:"content_path" json:"content_path" yaml:"content_path"`
	OutputPattern string `toml:"output_pattern" json:"output_pattern" yaml:"output_pattern"`
	Template      string `toml:"template" json:"template" yaml:"template"`
	// TODO: atom feed support
	GenerateRSS bool   `toml:"generate_rss" json:"generate_rss" yaml:"generate_rss"`
	Publish     string `toml:"publish" json:"publish" yaml:"publish"`
}

var defaultConfig = Config{
	PublishDir:       "public",
	StaticDir:        "static",
	ContentDir:       "content",
	TemplateDir:      "templates",
	TemplateExt:      ".tmpl",
	CacheDir:         "build",
	BaseTemplatePath: "",
	Copyright:        "",
	CleanBuild:       false,
	BuildDraft:       false,
	BuildFuture:      false,
	AllowUnsafeHTML:  false,
	RSS: RSS{
		Generate:    false,
		Title:       "Your Site",
		Description: "built with Cedar",
		URL:         "example.com",
	},
	ATProto: ATProto{},
	Routes: []Route{
		{
			ContentPath:   "index.md",
			OutputPattern: "/",
			Template:      "_index.html",
		},
	},
}

type decoder interface {
	Decode(v any) error
}

func Parse(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)

	var d decoder
	switch ext := filepath.Ext(path); ext {
	case ".toml":
		d = toml.NewDecoder(buf)
	case ".json":
		d = json.NewDecoder(buf)
	case ".yaml":
		d = yaml.NewDecoder(buf)
	default:
		return nil, fmt.Errorf("invalid config file type: %s", ext)
	}

	cfg := defaultConfig
	if err := d.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
