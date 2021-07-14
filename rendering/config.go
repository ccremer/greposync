package rendering

import (
	"os"
	"path"
	"path/filepath"
	"strings"

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
		r.p.WarnF("file not loaded: %s", err)
	}
	return nil
}

func (r *Renderer) loadDataForFile(fileName string) (Values, error) {
	// Load the global variables into exposed values
	data := make(Values)
	err := r.k.Unmarshal(":globals", &data)
	if err != nil {
		return data, err
	}
	segments := strings.Split(fileName, string(filepath.Separator))
	filePath := ""
	// Load the top-dir first (if any), then subdirs, then file-specific variables into values
	for _, segment := range segments {
		filePath = path.Join(filePath, segment)
		// Values applicable for directories are in the form of "my-dir/", otherwise they could be files.
		if filePath != fileName {
			filePath += string(filepath.Separator)
		}
		err = r.k.Unmarshal(filePath, &data)
		if err != nil {
			return data, err
		}
	}
	return data, err
}

func fileExists(fileName string) bool {
	if info, err := os.Stat(fileName); err == nil && !info.IsDir() {
		return true
	}
	return false
}
