package generator

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ptdewey/cedar/internal/atproto"
	"github.com/ptdewey/cedar/internal/config"
)

// WriteVerificationEndpoint generates the .well-known/site.standard.publication file
// containing the AT-URIs of all publication records, one per line.
func WriteVerificationEndpoint(cfg *config.Config) error {
	state, err := atproto.LoadPublishState()
	if err != nil || len(state.Publications) == 0 {
		return nil
	}

	var uris []string
	for _, pub := range state.Publications {
		if pub.ATURI != "" {
			uris = append(uris, pub.ATURI)
		}
	}
	if len(uris) == 0 {
		return nil
	}

	wellKnownDir := filepath.Join(cfg.PublishDir, ".well-known")
	if err := os.MkdirAll(wellKnownDir, 0755); err != nil {
		return err
	}

	return os.WriteFile(
		filepath.Join(wellKnownDir, "site.standard.publication"),
		[]byte(strings.Join(uris, "\n")),
		0644,
	)
}
