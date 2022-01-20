package cfg

import (
	"strings"

	"github.com/ccremer/greposync/cfg/flag"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/urfave/cli/v2"
)

// ParseConfig overrides given config defaults from file and with environment variables.
func ParseConfig(configPath string, config *Configuration, ctx *cli.Context) error {
	koanfInstance := koanf.New(".")

	// Load file
	if err := koanfInstance.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		return err
	}

	// Environment variables
	if err := koanfInstance.Load(env.Provider("", ".", func(s string) string {
		/*
			Configuration can contain hierarchies (YAML, etc.) and CLI flags dashes.
			To read environment variables with hierarchies and dashes we replace the hierarchy delimiter with double underscore and dashes with single underscore.
			So that parent.child-with-dash becomes PARENT__CHILD_WITH_DASH
		*/
		s = strings.Replace(strings.ToLower(s), "__", ".", -1)
		s = strings.Replace(strings.ToLower(s), "_", "-", -1)
		return s
	}), nil); err != nil {
		return err
	}

	// CLI flags
	if err := koanfInstance.Load(flag.Provider(ctx, "-", koanfInstance, nil), nil); err != nil {
		return err
	}

	return koanfInstance.Unmarshal("", &config)
}
