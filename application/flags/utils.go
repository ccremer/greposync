package flags

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func Prefixed(names ...string) []string {
	for i := 0; i < len(names); i++ {
		names[i] = EnvVarPrefix + names[i]
	}
	return names
}

func FromYAML(flags []cli.Flag) cli.BeforeFunc {
	return altsrc.InitInputSource(flags, func() (altsrc.InputSourceContext, error) {
		return altsrc.NewYamlSourceFromFile("greposync.yml")
	})
}

func And(fns ...func(ctx *cli.Context) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		for _, fn := range fns {
			if err := fn(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}

// CollectFlagValues returns a map with keys being the flag names and their values being flag values.
// The result contains the flags only of the current command.
// The result contains all flags, not just the ones that have been set.
func CollectFlagValues(ctx *cli.Context) map[string]interface{} {
	flags := ctx.Command.Flags
	res := make(map[string]interface{})
	for _, flag := range flags {
		flagName := flag.Names()[0]
		if flagName == "help" {
			continue
		}
		value := ctx.Value(flagName)
		if strValue, ok := value.(string); ok {
			if strValue == "" {
				value = "''"
			}
		}
		res[flagName] = value
	}
	return res
}
