package main

import (
	"fmt"
	"os"

	"github.com/ptdewey/cedar/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "build":
		err = cli.RunBuild(os.Args[2:])
	case "auth":
		err = cli.RunAuth(os.Args[2:])
	case "publish":
		err = cli.RunPublish(os.Args[2:])
	case "convert":
		err = cli.RunConvert(os.Args[2:])
	default:
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: cedar <command> [flags]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  build     Build the static site")
	fmt.Fprintln(os.Stderr, "  auth      Authenticate with ATProto via OAuth")
	fmt.Fprintln(os.Stderr, "  publish   Publish/sync content to ATProto PDS")
	fmt.Fprintln(os.Stderr, "  convert   Convert markdown files to Leaflet JSON or HTML")
}
