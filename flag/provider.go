package flag

import (
	"errors"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/maps"
	"github.com/urfave/cli/v2"
)

// Cli implements a urfave/cli command line provider.
type Cli struct {
	delim string
	ko    *koanf.Koanf
	ctx   *cli.Context
}

/*
Provider returns a commandline flags provider that returns a nested map[string]interface{} where the nesting hierarchy of keys are defined by delim.
For instance, the delim "." will convert the flag name `parent.child.key: 1` to `{parent: {child: {key: 1}}}`.
It takes an optional (but recommended) Koanf instance to see if the the flags defined have been set from other providers, for instance, a config file.
If they are not, then the default values of the flags are merged.
If they do exist, the flag values are not merged but only the values that have been explicitly set in the command line are merged.
*/
func Provider(ctx *cli.Context, delim string, ko *koanf.Koanf) *Cli {
	return &Cli{
		delim: delim,
		ko:    ko,
		ctx:   ctx,
	}
}

// Read reads the flag variables and returns a nested conf map.
func (p *Cli) Read() (map[string]interface{}, error) {
	mp := make(map[string]interface{})
	if p.ctx.Command == nil {
		return mp, nil
	}
	for _, flag := range p.ctx.Command.Flags {
		flagName := flag.Names()[0]
		val := p.ctx.Value(flagName)

		// If the default value of the flag was never changed by the user,
		// it should not override the value in the conf map (if it exists in the first place).
		if !p.ctx.IsSet(flagName) {
			if p.ko != nil {
				if p.ko.Exists(flagName) {
					continue
				}
			} else {
				continue
			}
		}
		mp[flagName] = val
	}
	return maps.Unflatten(mp, p.delim), nil
}

// ReadBytes is not supported.
func (p *Cli) ReadBytes() ([]byte, error) {
	return nil, errors.New("cli provider does not support this method")
}

// Watch is not supported.
func (p *Cli) Watch(cb func(event interface{}, err error)) error {
	return errors.New("cli provider does not support this method")
}
