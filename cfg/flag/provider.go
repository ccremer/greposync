package flag

import (
	"errors"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/maps"
	"github.com/urfave/cli/v2"
)

// Cli implements a urfave/cli command line provider.
type Cli struct {
	delim       string
	ko          *koanf.Koanf
	ctx         *cli.Context
	aliases     map[string]string
	useDefaults bool
}

/*
Provider returns a commandline flags provider that returns a nested map[string]interface{} where the nesting hierarchy of keys is defined by flagDelim.
For instance, the flagDelim "." will convert the flag name `parent.child.key: 1` to `{parent: {child: {key: 1}}}`.
It takes an optional (but recommended) Koanf instance to see if the flags defined have been set from other providers, for instance, a config file.
If there are pre-existing values, then they are overwritten.
The aliases map allows putting flag values into trees at arbitrary positions.
For example, given alias["flag"] = "nested-flag", flagDelim = "-" and parsed arguments = "--flag=bar", then the resulting map is `{nested: {flag: "bar"}}`.
*/
func Provider(ctx *cli.Context, flagDelim string, ko *koanf.Koanf, aliases map[string]string) *Cli {
	if aliases == nil {
		aliases = map[string]string{}
	}
	return &Cli{
		delim:   flagDelim,
		ko:      ko,
		ctx:     ctx,
		aliases: aliases,
	}
}

// Read reads the flag variables and returns a nested conf map.
func (p *Cli) Read() (map[string]interface{}, error) {
	mp := make(map[string]interface{})
	for _, flagName := range p.ctx.FlagNames() {
		val := p.ctx.Value(flagName)

		alias := p.aliases[flagName]
		if alias == "" {
			alias = flagName
		}

		if p.ctx.IsSet(flagName) {
			mp[alias] = val
		}
	}
	return maps.Unflatten(mp, p.delim), nil
}

// ReadBytes is not supported.
func (p *Cli) ReadBytes() ([]byte, error) {
	return nil, errors.New("cli provider does not support this method")
}

// Watch is not supported.
func (p *Cli) Watch(_ func(event interface{}, err error)) error {
	return errors.New("cli provider does not support this method")
}
