package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	libleaflet "github.com/ptdewey/standard-site-go/leaflet"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

func RunConvert(args []string) error {
	fs := flag.NewFlagSet("convert", flag.ExitOnError)
	html := fs.Bool("html", false, "output HTML preview instead of JSON")
	outDir := fs.String("out", "", "write output files to this directory instead of stdout")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: cedar convert [flags] <file.md> [file.md ...]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Convert markdown files to Leaflet records.")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Flags:")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	files := fs.Args()
	if len(files) == 0 {
		fs.Usage()
		return fmt.Errorf("no input files specified")
	}

	if *outDir != "" {
		if err := os.MkdirAll(*outDir, 0o755); err != nil {
			return fmt.Errorf("creating output directory: %w", err)
		}
	}

	md := goldmark.New(goldmark.WithExtensions(extension.GFM, extension.Footnote))

	for _, path := range files {
		if err := convertFile(md, path, *html, *outDir); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
	}
	return nil
}

func convertFile(md goldmark.Markdown, path string, asHTML bool, outDir string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	source := stripFrontmatter(raw)
	doc := md.Parser().Parse(text.NewReader(source))
	content := libleaflet.Convert(source, doc)

	if asHTML {
		body := libleaflet.RenderHTML(content)
		out := "<!DOCTYPE html>\n<html><head><meta charset=\"utf-8\"></head><body>\n" + body + "\n</body></html>\n"
		return writeOutput([]byte(out), path, ".html", outDir)
	}

	out, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return err
	}
	return writeOutput(append(out, '\n'), path, ".json", outDir)
}

// writeOutput writes data to stdout or, if outDir is set, to a file named
// after the input path with the given extension.
func writeOutput(data []byte, inputPath, ext, outDir string) error {
	if outDir == "" {
		_, err := os.Stdout.Write(data)
		return err
	}

	base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	outPath := filepath.Join(outDir, base+ext)
	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", outPath)
	return nil
}

// stripFrontmatter removes YAML/TOML front matter delimited by "---" lines.
func stripFrontmatter(src []byte) []byte {
	s := string(src)
	if !strings.HasPrefix(s, "---") {
		return src
	}
	parts := strings.SplitN(s, "---", 3)
	if len(parts) < 3 {
		return src
	}
	return []byte(strings.TrimSpace(parts[2]))
}
