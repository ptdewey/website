package cli

import (
	"flag"

	"github.com/ptdewey/cedar/internal/config"
)

func addConfigFlag(fs *flag.FlagSet) *string {
	return fs.String("config", "cedar.toml", "path to config file")
}

func parseConfig(path string) (*config.Config, error) {
	return config.Parse(path)
}
