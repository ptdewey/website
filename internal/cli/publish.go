package cli

import (
	"flag"
	"fmt"

	"github.com/ptdewey/cedar/internal/atproto"
	"github.com/ptdewey/cedar/internal/parser"
)

// previewDir is the default directory for leaflet preview HTML files.
const previewDir = "_preview"

func RunPublish(args []string) error {
	fs := flag.NewFlagSet("publish", flag.ExitOnError)
	configPath := addConfigFlag(fs)
	dryRun := fs.Bool("dry-run", false, "print records as JSON without publishing")
	preview := fs.Bool("preview", false, "generate leaflet HTML preview files in _preview/")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := parseConfig(*configPath)
	if err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}

	if cfg.ATProto.Handle == "" {
		return fmt.Errorf("atproto.handle must be set in config")
	}
	if len(cfg.ATProto.Publications) == 0 {
		return fmt.Errorf("at least one publication must be configured under [atproto.publications]")
	}

	pages, err := parser.ProcessDirectory(cfg)
	if err != nil {
		return fmt.Errorf("failed to process content directory: %w", err)
	}

	if *dryRun {
		return atproto.DryRun(cfg, pages)
	}

	if *preview {
		return atproto.Preview(cfg, pages, previewDir)
	}

	return atproto.Publish(cfg, pages)
}
