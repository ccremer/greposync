package rendering

import (
	"github.com/ccremer/git-repo-sync/printer"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var k = koanf.New(".")

func LoadConfigFile(path string) {
	// Load YAML config and merge into the previously loaded config (because we can).
	err := k.Load(file.Provider(path), yaml.Parser())
	printer.CheckIfError(err)
}

func Unmarshal(file string) map[string]interface{} {

	m := make(map[string]interface{})

	err := k.Unmarshal(":globals", &m)
	err = k.Unmarshal(file, &m)
	printer.CheckIfError(err)
	return m
}
