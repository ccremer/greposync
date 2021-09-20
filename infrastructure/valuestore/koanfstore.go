package valuestore

import (
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ccremer/greposync/domain"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

// KoanfValueStore implements domain.ValueStore.
type KoanfValueStore struct {
	m               *sync.Mutex
	globalKoanf     *koanf.Koanf
	instrumentation *ValueStoreInstrumentation
	cache           map[*domain.GitURL]*koanf.Koanf
}

var (
	SyncConfigFileName   = ".sync.yml"
	GlobalConfigFileName = "config_defaults.yml"
)

// NewValueStore returns a new instance of domain.ValueStore.
func NewValueStore(instrumentation *ValueStoreInstrumentation) *KoanfValueStore {
	return &KoanfValueStore{
		instrumentation: instrumentation,
		cache:           map[*domain.GitURL]*koanf.Koanf{},
		m:               &sync.Mutex{},
	}
}

// FetchValuesForTemplate implements domain.ValueStore.
func (k *KoanfValueStore) FetchValuesForTemplate(template *domain.Template, config *domain.GitRepository) (domain.Values, error) {
	k.loadGlobals()
	repoKoanf, err := k.prepareRepoKoanf(config)
	if err != nil {
		return domain.Values{}, err
	}
	return k.loadDataForTemplate(repoKoanf, template.RelativePath.String())
}

// FetchUnmanagedFlag implements domain.ValueStore.
func (k *KoanfValueStore) FetchUnmanagedFlag(template *domain.Template, config *domain.GitRepository) (bool, error) {
	k.loadGlobals()
	repoKoanf, err := k.prepareRepoKoanf(config)
	if err != nil {
		return false, err
	}
	return k.loadBooleanFlag(repoKoanf, template.RelativePath.String(), "unmanaged")
}

// FetchTargetPath implements domain.ValueStore.
func (k *KoanfValueStore) FetchTargetPath(template *domain.Template, config *domain.GitRepository) (domain.Path, error) {
	k.loadGlobals()
	repoKoanf, err := k.prepareRepoKoanf(config)
	if err != nil {
		return "", err
	}
	return k.loadTargetPath(repoKoanf, template.RelativePath.String())
}

// FetchFilesToDelete implements domain.ValueStore.
func (k *KoanfValueStore) FetchFilesToDelete(config *domain.GitRepository) ([]domain.Path, error) {
	k.loadGlobals()
	repoKoanf, err := k.prepareRepoKoanf(config)
	if err != nil {
		return []domain.Path{}, err
	}
	return k.loadFilesToDelete(repoKoanf)
}

func (k *KoanfValueStore) prepareRepoKoanf(repository *domain.GitRepository) (*koanf.Koanf, error) {
	if repoKoanf, exists := k.cache[repository.URL]; exists {
		return repoKoanf, nil
	}
	repoKoanf, err := k.loadAndMergeConfig(repository)
	if err != nil {
		return nil, err
	}
	k.cache[repository.URL] = repoKoanf
	return repoKoanf, nil
}

func (k *KoanfValueStore) loadAndMergeConfig(repository *domain.GitRepository) (*koanf.Koanf, error) {
	syncFile := path.Join(repository.RootDir.String(), SyncConfigFileName)
	repoKoanf := koanf.New(".")
	// Load the config from config_defaults.yml
	err := repoKoanf.Merge(k.globalKoanf)
	if err != nil {
		return repoKoanf, err
	}
	k.instrumentation.attemptingLoadConfig(repository.URL.GetFullName(), syncFile)
	// Load the config from .sync.yml
	err = repoKoanf.Load(file.Provider(syncFile), yaml.Parser())
	return repoKoanf, k.instrumentation.loadedConfigIfNil(repository.URL.GetFullName(), err)
}

func (k *KoanfValueStore) loadDataForTemplate(repoKoanf *koanf.Koanf, templateFileName string) (domain.Values, error) {
	// Load the global variables into exposed values
	data := make(domain.Values)
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
	k.instrumentation.attemptingLoadConfig("", GlobalConfigFileName)
	// Load the config from global config file, but ignore errors.
	err := k.globalKoanf.Load(file.Provider(GlobalConfigFileName), yaml.Parser())
	_ = k.instrumentation.loadedConfigIfNil("", err)
}
