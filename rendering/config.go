package rendering

import (
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

func (r *Renderer) loadVariables(syncFile string) error {
	// Load the config from config_defaults.yml
	err := r.k.Merge(r.globalDefaults)
	if err != nil {
		return err
	}
	r.p.DebugF("Loading sync config %s", syncFile)
	// Load the config from .sync.yml
	err = r.k.Load(file.Provider(syncFile), yaml.Parser())
	if err != nil {
		r.p.WarnF("file %s not loaded: %s", syncFile, err)
	}
	return nil
}

func (r *Renderer) loadDataForFile(fileName string) (Values, error) {
	// Load the global variables into exposed values
	data := make(Values)
	err := r.k.Unmarshal(":globals", &data)
	if err != nil {
		return nil, err
	}
	// Load the file-specific variables into values
	err = r.k.Unmarshal(fileName, &data)
	return data, err
}

func (r *Renderer) fileExists(fileName string) bool {
	if info, err := os.Stat(fileName); err != nil || info.IsDir() {
		return false
	}
	return true
}
