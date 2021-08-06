package valuestore

import (
	"errors"
	"path"
	"path/filepath"
	"strings"

	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/printer"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

// KoanfValueStore implements core.ValueStore.
type KoanfValueStore struct {
	globalKoanf *koanf.Koanf
	log         printer.Printer
	cache       map[string]*koanf.Koanf
}

var (
	SyncConfigFileName = ".sync.yml"
)

// NewValueStore returns a new instance of core.ValueStore.
func NewValueStore(globalKoanf *koanf.Koanf) *KoanfValueStore {
	return &KoanfValueStore{
		globalKoanf: globalKoanf,
		log:         printer.New(),
		cache:       map[string]*koanf.Koanf{},
	}
}

// FetchValuesForTemplate implements core.ValueStore.
func (k *KoanfValueStore) FetchValuesForTemplate(template core.Template, config *core.GitRepositoryConfig) (core.Values, error) {
	repoKoanf, err := k.prepareRepoKoanf(config)
	if err != nil {
		return core.Values{}, err
	}
	return k.loadDataForTemplate(repoKoanf, template.GetRelativePath())
}

// FetchUnmanagedFlag implements core.ValueStore.
func (k *KoanfValueStore) FetchUnmanagedFlag(template core.Template, config *core.GitRepositoryConfig) (bool, error) {
	repoKoanf, err := k.prepareRepoKoanf(config)
	if err != nil {
		return false, err
	}
	return k.loadBooleanFlag(repoKoanf, template.GetRelativePath(), "unmanaged")
}

// FetchTargetPath implements core.ValueStore.
func (k *KoanfValueStore) FetchTargetPath(template core.Template, config *core.GitRepositoryConfig) (string, error) {
	repoKoanf, err := k.prepareRepoKoanf(config)
	if err != nil {
		return "", err
	}
	return k.loadTargetPath(repoKoanf, template.GetRelativePath())
}

// FetchFilesToDelete implements core.ValueStore.
func (k *KoanfValueStore) FetchFilesToDelete(config *core.GitRepositoryConfig) ([]string, error) {
	repoKoanf, err := k.prepareRepoKoanf(config)
	if err != nil {
		return []string{}, err
	}
	return k.loadFilesToDelete(repoKoanf)
}

func (k *KoanfValueStore) prepareRepoKoanf(config *core.GitRepositoryConfig) (*koanf.Koanf, error) {
	k.log.SetName(config.URL.GetRepositoryName())
	if repoKoanf, exists := k.cache[config.RootDir]; exists {
		return repoKoanf, nil
	}
	repoKoanf, err := k.loadAndMergeConfig(path.Join(config.RootDir, SyncConfigFileName))
	if err != nil {
		return nil, err
	}
	k.cache[config.RootDir] = repoKoanf
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

func (k *KoanfValueStore) loadFilesToDelete(repoKoanf *koanf.Koanf) ([]string, error) {
	filePaths := make([]string, 0)
	allKeys := repoKoanf.Raw()
	// Go through all top-level keys, which are the file names
	for filePath, _ := range allKeys {
		// If the filename is already handled by the template renderer, ignore it.
		// Otherwise, add files that have deletion flag, but ignore directories
		if !pathIsFile(filePath) {
			continue
		}
		if filePath == ":globals" {
			// can't delete file named ':globals' anyway
			continue
		}
		del, err := k.loadBooleanFlag(repoKoanf, filePath, "delete")
		if errors.Is(err, core.ErrKeyNotFound) {
			continue
		}
		if del {
			filePaths = append(filePaths, filePath)
		}
	}
	return filePaths, nil
}

func (k *KoanfValueStore) loadBooleanFlag(repoKoanf *koanf.Koanf, relativePath, flagName string) (bool, error) {
	values, err := k.loadDataForTemplate(repoKoanf, relativePath)
	if err != nil {
		return false, err
	}
	flag, exists := values[flagName]
	if exists {
		return flag == true, nil
	}
	return false, core.ErrKeyNotFound
}

func (k *KoanfValueStore) loadTargetPath(repoKoanf *koanf.Koanf, relativePath string) (string, error) {
	values, err := k.loadDataForTemplate(repoKoanf, relativePath)
	if err != nil {
		return "", err
	}
	targetPath, exists := values["targetPath"]
	if exists {
		newPath, isString := targetPath.(string)
		if isString {
			if strings.HasSuffix(newPath, "/") {
				return path.Clean(path.Join(newPath, path.Base(relativePath))), nil
			}
			return newPath, nil
		}
		return "", nil
	}
	return "", nil
}

func pathIsFile(filePath string) bool {
	return !strings.HasSuffix(filePath, "/")
}
