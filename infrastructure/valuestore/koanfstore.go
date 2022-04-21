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

// KoanfStore implements domain.ValueStore.
type KoanfStore struct {
	m                  *sync.Mutex
	globalKoanf        *koanf.Koanf
	instrumentation    *ValueStoreInstrumentation
	cache              map[*domain.GitURL]*koanf.Koanf
	syncConfigFileName string
}

// NewKoanfStore returns a new instance of domain.ValueStore.
func NewKoanfStore(instrumentation *ValueStoreInstrumentation) *KoanfStore {
	return &KoanfStore{
		instrumentation:    instrumentation,
		cache:              map[*domain.GitURL]*koanf.Koanf{},
		m:                  &sync.Mutex{},
		syncConfigFileName: SyncConfigFileName,
	}
}

// FetchValuesForTemplate implements domain.ValueStore.
func (s *KoanfStore) FetchValuesForTemplate(template *domain.Template, config *domain.GitRepository) (domain.Values, error) {
	s.loadGlobals()
	repoKoanf, err := s.prepareRepoKoanf(config)
	if err != nil {
		return domain.Values{}, err
	}
	return s.loadValuesForTemplate(repoKoanf, template.CleanPath().String())
}

// FetchUnmanagedFlag implements domain.ValueStore.
func (s *KoanfStore) FetchUnmanagedFlag(template *domain.Template, config *domain.GitRepository) (bool, error) {
	s.loadGlobals()
	repoKoanf, err := s.prepareRepoKoanf(config)
	if err != nil {
		return false, err
	}
	return s.loadBooleanFlag(repoKoanf, template.CleanPath().String(), "unmanaged")
}

// FetchTargetPath implements domain.ValueStore.
func (s *KoanfStore) FetchTargetPath(template *domain.Template, config *domain.GitRepository) (domain.Path, error) {
	s.loadGlobals()
	repoKoanf, err := s.prepareRepoKoanf(config)
	if err != nil {
		return "", err
	}
	return s.loadTargetPath(repoKoanf, template.CleanPath().String())
}

// FetchFilesToDelete implements domain.ValueStore.
func (s *KoanfStore) FetchFilesToDelete(config *domain.GitRepository) ([]domain.Path, error) {
	s.loadGlobals()
	repoKoanf, err := s.prepareRepoKoanf(config)
	if err != nil {
		return []domain.Path{}, err
	}
	return s.loadFilesToDelete(repoKoanf)
}

func (s *KoanfStore) prepareRepoKoanf(repository *domain.GitRepository) (*koanf.Koanf, error) {
	if repoKoanf, exists := s.cache[repository.URL]; exists {
		return repoKoanf, nil
	}
	repoKoanf, err := s.loadAndMergeConfig(repository)
	if err != nil {
		return nil, err
	}
	s.cache[repository.URL] = repoKoanf
	return repoKoanf, nil
}

func (s *KoanfStore) loadValuesForTemplate(repoKoanf *koanf.Koanf, templateFileName string) (domain.Values, error) {
	// Load the global variables into exposed values
	merged := repoKoanf.Cut(":globals")
	segments := strings.Split(templateFileName, string(filepath.Separator))
	filePath := ""
	// Load the top-dir first (if any), then subdirs, then file-specific variables into values
	for _, segment := range segments {
		filePath = path.Join(filePath, segment)
		// Values applicable for directories are in the form of "my-dir/", otherwise they could be files.
		if filePath != templateFileName {
			filePath += string(filepath.Separator)
		}
		err := merged.Merge(repoKoanf.Cut(filePath))
		if err != nil {
			return nil, err
		}
	}
	return merged.Raw(), nil
}

func (s *KoanfStore) loadAndMergeConfig(repository *domain.GitRepository) (*koanf.Koanf, error) {
	syncFile := path.Join(repository.RootDir.String(), s.syncConfigFileName)
	repoKoanf := koanf.New("")
	// Load the config from config_defaults.yml
	err := repoKoanf.Merge(s.globalKoanf)
	if err != nil {
		return repoKoanf, err
	}
	s.instrumentation.attemptingLoadConfig(repository.URL.GetFullName(), syncFile)
	err = s.loadYaml(repoKoanf, syncFile)
	return repoKoanf, s.instrumentation.loadedConfigIfNil(repository.URL.GetFullName(), err)
}

func (s *KoanfStore) loadYaml(repoKoanf *koanf.Koanf, syncFilePath string) error {
	// Load the config from .sync.yml
	err := repoKoanf.Load(file.Provider(syncFilePath), yaml.Parser())
	return err
}

func (s *KoanfStore) loadGlobals() {
	s.m.Lock()
	defer s.m.Unlock()
	if s.globalKoanf != nil {
		return
	}

	s.globalKoanf = koanf.New(".")
	s.instrumentation.attemptingLoadConfig("", GlobalConfigFileName)
	// Load the config from global config file, but ignore errors.
	err := s.globalKoanf.Load(file.Provider(GlobalConfigFileName), yaml.Parser())
	_ = s.instrumentation.loadedConfigIfNil("", err)
}
