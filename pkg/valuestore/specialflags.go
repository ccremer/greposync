package valuestore

import (
	"errors"
	"path"
	"strings"

	"github.com/ccremer/greposync/core"
	"github.com/knadh/koanf"
)

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

func pathIsFile(filePath string) bool {
	return !strings.HasSuffix(filePath, "/")
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
