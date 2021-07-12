package cfg

import (
	"strings"

	"github.com/ccremer/greposync/cfg/cli"
	"github.com/ccremer/greposync/printer"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	urfavecli "github.com/urfave/cli/v2"
)

// ParseConfig overrides given config defaults from file and with environment variables.
func ParseConfig(configPath string, config *Configuration, ctx *urfavecli.Context) error {
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
	if err := koanfInstance.Load(cli.Provider(ctx, "-", koanfInstance), nil); err != nil {
		return err
	}

	return koanfInstance.Unmarshal("", &config)
}

// Sanitize does corrective actions on the configuration hierarchy.
func (config *Configuration) Sanitize() {
	level, err := printer.ParseLogLevel(config.Log.Level)
	if err != nil {
		printer.WarnF("Could not parse log level, fallback to default level")
	}
	printer.DefaultPrinter.SetLevel(level)
}
