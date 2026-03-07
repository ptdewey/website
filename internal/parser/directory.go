package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ptdewey/cedar/internal/config"
)

func ProcessDirectory(cfg *config.Config) ([]Page, error) {
	var pages []Page

	for _, route := range cfg.Routes {
		// TODO: allow creating routes from templates w/o content file

		contentPath := filepath.Join(cfg.ContentDir, route.ContentPath)

		info, err := os.Stat(contentPath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("route '%s' points to non-existent path: %s", route.ContentPath, contentPath)
			}

			return nil, fmt.Errorf("error accessing route '%s': %w", route.ContentPath, err)
		}

		if info.IsDir() {
			dirPages, err := processRouteDirectory(contentPath, cfg)
			if err != nil {
				return nil, err
			}
			if len(dirPages) == 0 {
				return nil, fmt.Errorf("route '%s' is a directory but contains no markdown files", route.ContentPath)
			}
			pages = append(pages, dirPages...)
		} else {
			if !strings.HasSuffix(contentPath, ".md") {
				return nil, fmt.Errorf("route '%s' points to non-markdown file: %s", route.ContentPath, contentPath)
			}

			page, err := ProcessMarkdownFile(contentPath, cfg)
			if err != nil {
				return nil, fmt.Errorf("error processing route '%s': %w", route.ContentPath, err)
			}

			pages = append(pages, page)
		}
	}

	return pages, nil
}

func processRouteDirectory(dir string, cfg *config.Config) ([]Page, error) {
	var pages []Page

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), ".md") {
			page, err := ProcessMarkdownFile(path, cfg)
			if err != nil {
				return err
			}
			pages = append(pages, page)
		}
		return nil
	})

	return pages, err
}
