package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/ptdewey/cedar/internal/generator"
	"github.com/ptdewey/cedar/internal/parser"
	"github.com/ptdewey/cedar/internal/rss"
)

func RunBuild(args []string) error {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	configPath := addConfigFlag(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := parseConfig(*configPath)
	if err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}

	pages, err := parser.ProcessDirectory(cfg)
	if err != nil {
		return fmt.Errorf("failed to process content directory: %w", err)
	}

	if err := generator.WriteHTMLFiles(pages, cfg.CacheDir, cfg); err != nil {
		return fmt.Errorf("error writing HTML files: %w", err)
	}

	if cfg.CleanBuild {
		if err := os.RemoveAll(cfg.PublishDir); err != nil {
			fmt.Printf("failed to remove publish directory '%s': %v\n", cfg.PublishDir, err)
		}
	}

	if err := generator.CopyDirIncremental(cfg.CacheDir, cfg.PublishDir); err != nil {
		return fmt.Errorf("error copying build cache: %w", err)
	}
	_ = os.RemoveAll(cfg.CacheDir)

	if err := generator.CopyDirIncremental(cfg.StaticDir, cfg.PublishDir); err != nil {
		return fmt.Errorf("error copying static directory: %w", err)
	}

	if cfg.RSS.Generate {
		if err := rss.GenerateRSS(pages, cfg); err != nil {
			return fmt.Errorf("error writing rss.xml: %w", err)
		}
	}

	if cfg.ATProto.Handle != "" {
		if err := generator.WriteVerificationEndpoint(cfg); err != nil {
			return fmt.Errorf("error writing verification endpoint: %w", err)
		}
	}

	msg := "Successfully generated HTML files"
	if cfg.RSS.Generate {
		msg += " and rss.xml"
	}
	fmt.Println(msg)
	return nil
}
