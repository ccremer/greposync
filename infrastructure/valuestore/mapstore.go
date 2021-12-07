package valuestore

import (
	"os"
	"path"
	"strings"
	"sync"

	"github.com/ccremer/greposync/domain"
	"github.com/imdario/mergo"
	"sigs.k8s.io/yaml"
)

// MapStore implements domain.ValueStore using map structure as backend.
// It comes with per-repository transparent caching and lazy-loading.
// It also features a global config file that is loaded on first access.
type MapStore struct {
	m               *sync.Mutex
	globalConfig    config
	instrumentation *ValueStoreInstrumentation
	cache           map[*domain.GitURL]config
}

// config is just an alias for easier readability.
type config map[string]interface{}

// NewMapStore returns a new instance of domain.ValueStore.
func NewMapStore(instrumentation *ValueStoreInstrumentation) *MapStore {
	return &MapStore{
		instrumentation: instrumentation,
		cache:           map[*domain.GitURL]config{},
		m:               &sync.Mutex{},
	}
}

// FetchValuesForTemplate implements domain.ValueStore.
func (s *MapStore) FetchValuesForTemplate(template *domain.Template, repository *domain.GitRepository) (domain.Values, error) {
	s.loadGlobals()
	repoConfig, err := s.prepareRepoConfig(repository)
	if err != nil {
		return domain.Values{}, err
	}
	return s.loadValuesForTemplate(repoConfig, template.CleanPath().String())
}

// FetchUnmanagedFlag implements domain.ValueStore.
func (s *MapStore) FetchUnmanagedFlag(template *domain.Template, repository *domain.GitRepository) (bool, error) {
	s.loadGlobals()
	repoConfig, err := s.prepareRepoConfig(repository)
	if err != nil {
		return false, err
	}
	return s.loadBooleanFlag(repoConfig, template.CleanPath().String(), "unmanaged")
}

// FetchTargetPath implements domain.ValueStore.
func (s *MapStore) FetchTargetPath(template *domain.Template, repository *domain.GitRepository) (domain.Path, error) {
	s.loadGlobals()
	repoConfig, err := s.prepareRepoConfig(repository)
	if err != nil {
		return "", err
	}
	return s.loadTargetPath(repoConfig, template.CleanPath().String())
}

// FetchFilesToDelete implements domain.ValueStore.
func (s *MapStore) FetchFilesToDelete(repository *domain.GitRepository) ([]domain.Path, error) {
	s.loadGlobals()
	repoConfig, err := s.prepareRepoConfig(repository)
	if err != nil {
		return []domain.Path{}, err
	}
	return s.loadFilesToDelete(repoConfig)
}

func (s *MapStore) prepareRepoConfig(repository *domain.GitRepository) (config, error) {
	if cfg, exists := s.cache[repository.URL]; exists {
		return cfg, nil
	}
	repoConfig, err := s.loadAndMergeConfig(repository)
	if err != nil {
		return nil, err
	}
	s.cache[repository.URL] = repoConfig
	return repoConfig, nil
}

func (s *MapStore) loadAndMergeConfig(repository *domain.GitRepository) (config, error) {
	syncFile := path.Join(repository.RootDir.String(), SyncConfigFileName)
	repoConfig := config{}
	// Load the config from config_defaults.yml
	err := mergo.Merge(&repoConfig, s.globalConfig, mergo.WithOverride)
	if err != nil {
		return repoConfig, err
	}
	s.instrumentation.attemptingLoadConfig(repository.URL.GetFullName(), syncFile)
	// Load the config from .sync.yml
	syncCfg, err := s.loadYaml(syncFile)
	if err == nil {
		err = mergo.Merge(&repoConfig, syncCfg, mergo.WithOverride)
	}
	return repoConfig, s.instrumentation.loadedConfigIfNil(repository.URL.GetFullName(), err)
}

func (s *MapStore) loadValuesForTemplate(repoConfig config, templateFileName string) (domain.Values, error) {
	// Load the global variables into exposed values
	data := make(map[string]interface{})
	if globals, found := repoConfig[":globals"]; found {
		err := mergo.Merge(&data, globals, mergo.WithOverride)
		if err != nil {
			return data, err
		}
	}

	segments := strings.Split(templateFileName, "/")
	filePath := ""
	// Load the top-dir first (if any), then subdirs, then file-specific variables into values
	for _, segment := range segments {
		filePath = path.Join(filePath, segment)
		// Values applicable for directories are in the form of "my-dir/", otherwise they could be files.
		if filePath != templateFileName {
			filePath += "/"
		}
		if cfg, found := repoConfig[filePath]; found {
			err := mergo.Merge(&data, cfg, mergo.WithOverride)
			if err != nil {
				return data, err
			}
		}
	}
	return data, nil
}

func (s *MapStore) loadGlobals() {
	s.m.Lock()
	defer s.m.Unlock()
	if s.globalConfig != nil {
		return
	}

	s.instrumentation.attemptingLoadConfig("", GlobalConfigFileName)
	// Load the config from global config file, but ignore errors.
	c, err := s.loadYaml(GlobalConfigFileName)
	s.globalConfig = c
	_ = s.instrumentation.loadedConfigIfNil("", err)
}

func (s *MapStore) loadYaml(file string) (config, error) {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	c := config{}
	err = yaml.Unmarshal(yamlFile, &c)
	return c, err
}
