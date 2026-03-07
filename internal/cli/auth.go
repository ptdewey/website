package cli

import (
	"flag"
	"fmt"

	"github.com/ptdewey/cedar/internal/atproto"
)

func RunAuth(args []string) error {
	fs := flag.NewFlagSet("auth", flag.ExitOnError)
	configPath := addConfigFlag(fs)
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

	return atproto.RunOAuthFlow(cfg)
}
