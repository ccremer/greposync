package valuestore

import (
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/printer"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

// KoanfValueStore implements core.ValueStore.
type KoanfValueStore struct {
	m           *sync.Mutex
	globalKoanf *koanf.Koanf
	log         printer.Printer
	cache       map[*core.GitURL]*koanf.Koanf
}

var (
	SyncConfigFileName   = ".sync.yml"
	GlobalConfigFileName = "config_defaults.yml"
)

// NewValueStore returns a new instance of core.ValueStore.
func NewValueStore() *KoanfValueStore {
	return &KoanfValueStore{
		log:   printer.New(),
		cache: map[*core.GitURL]*koanf.Koanf{},
		m:     &sync.Mutex{},
	}
}

// FetchValuesForTemplate implements core.ValueStore.
func (k *KoanfValueStore) FetchValuesForTemplate(template core.Template, config *core.GitRepositoryProperties) (core.Values, error) {
	k.loadGlobals()
	repoKoanf, err := k.prepareRepoKoanf(config)
	if err != nil {
		return core.Values{}, err
	}
	return k.loadDataForTemplate(repoKoanf, template.GetRelativePath())
}

// FetchUnmanagedFlag implements core.ValueStore.
func (k *KoanfValueStore) FetchUnmanagedFlag(template core.Template, config *core.GitRepositoryProperties) (bool, error) {
	k.loadGlobals()
	repoKoanf, err := k.prepareRepoKoanf(config)
	if err != nil {
		return false, err
	}
	return k.loadBooleanFlag(repoKoanf, template.GetRelativePath(), "unmanaged")
}

// FetchTargetPath implements core.ValueStore.
func (k *KoanfValueStore) FetchTargetPath(template core.Template, config *core.GitRepositoryProperties) (string, error) {
	k.loadGlobals()
	repoKoanf, err := k.prepareRepoKoanf(config)
	if err != nil {
		return "", err
	}
	return k.loadTargetPath(repoKoanf, template.GetRelativePath())
}

// FetchFilesToDelete implements core.ValueStore.
func (k *KoanfValueStore) FetchFilesToDelete(config *core.GitRepositoryProperties) ([]string, error) {
	k.loadGlobals()
	repoKoanf, err := k.prepareRepoKoanf(config)
	if err != nil {
		return []string{}, err
	}
	return k.loadFilesToDelete(repoKoanf)
}

func (k *KoanfValueStore) prepareRepoKoanf(config *core.GitRepositoryProperties) (*koanf.Koanf, error) {
	k.log.SetName(config.URL.GetRepositoryName())
	if repoKoanf, exists := k.cache[config.URL]; exists {
		return repoKoanf, nil
	}
	repoKoanf, err := k.loadAndMergeConfig(path.Join(config.RootDir, SyncConfigFileName))
	if err != nil {
		return nil, err
	}
	k.cache[config.URL] = repoKoanf
	return repoKoanf, nil
}

func (k *KoanfValueStore) loadAndMergeConfig(syncFile string) (*koanf.Koanf, error) {
	repoKoanf := koanf.New(".")
	// Load the config from config_defaults.yml
	err := repoKoanf.Merge(k.globalKoanf)
	if err != nil {
		return repoKoanf, err
	}
	k.log.DebugF("Loading sync config '%s'", syncFile)
	// Load the config from .sync.yml
	err = repoKoanf.Load(file.Provider(syncFile), yaml.Parser())
	if err != nil {
		k.log.WarnF("file '%s' not loaded: %s", syncFile, err)
	}
	return repoKoanf, nil
}

func (k *KoanfValueStore) loadDataForTemplate(repoKoanf *koanf.Koanf, templateFileName string) (core.Values, error) {
	// Load the global variables into exposed values
	data := make(core.Values)
	err := repoKoanf.Unmarshal(":globals", &data)
	if err != nil {
		return data, err
	}
	segments := strings.Split(templateFileName, string(filepath.Separator))
	filePath := ""
	// Load the top-dir first (if any), then subdirs, then file-specific variables into values
	for _, segment := range segments {
		filePath = path.Join(filePath, segment)
		// Values applicable for directories are in the form of "my-dir/", otherwise they could be files.
		if filePath != templateFileName {
			filePath += string(filepath.Separator)
		}
		err = repoKoanf.Unmarshal(filePath, &data)
		if err != nil {
			return data, err
		}
	}
	return data, err
}

func (k *KoanfValueStore) loadGlobals() {
	k.m.Lock()
	defer k.m.Unlock()
	if k.globalKoanf != nil {
		return
	}
	k.globalKoanf = koanf.New(".")
	k.log.DebugF("Loading config '%s'", GlobalConfigFileName)
	// Load the config from global config file, but ignore errors.
	err := k.globalKoanf.Load(file.Provider(GlobalConfigFileName), yaml.Parser())
	if err != nil {
		k.log.WarnF("file '%s' not loaded: %s", GlobalConfigFileName, err)
	}
}
