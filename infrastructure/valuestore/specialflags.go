package valuestore

import (
	"errors"
	"path"
	"strings"

	"github.com/ccremer/greposync/domain"
	"github.com/knadh/koanf"
)

func (s *KoanfStore) loadFilesToDelete(repoConfig *koanf.Koanf, templates []*domain.Template) ([]domain.Path, error) {
	filePaths := make([]domain.Path, 0)
	// Go through all top-level keys, which are the file names.
	// We do this because the list of templates might not contain the desired file name to delete anymore.
	for filePath, _ := range repoConfig.Raw() {
		// If the filename is already handled by the template renderer, ignore it.
		// Otherwise, add files that have deletion flag, but ignore directories
		if !pathIsFile(filePath) {
			continue
		}
		if filePath == ":globals" {
			// can't delete file named ':globals' anyway
			continue
		}
		del, err := s.loadBooleanFlag(repoConfig, filePath, "delete")
		if errors.Is(err, domain.ErrKeyNotFound) {
			continue
		}
		p := domain.Path(filePath)
		if del && !p.IsInSlice(filePaths) {
			filePaths = append(filePaths, p)
		}
	}
	for _, template := range templates {
		del, err := s.loadBooleanFlag(repoConfig, template.CleanPath().String(), "delete")
		if errors.Is(err, domain.ErrKeyNotFound) {
			continue
		}
		if del && !template.RelativePath.IsInSlice(filePaths) {
			filePaths = append(filePaths, template.RelativePath)
		}
	}
	return filePaths, nil
}

func pathIsFile(filePath string) bool {
	return !strings.HasSuffix(filePath, "/")
}

func (s *KoanfStore) loadBooleanFlag(repoConfig *koanf.Koanf, relativePath, flagName string) (bool, error) {
	values, err := s.loadValuesForTemplate(repoConfig, relativePath)
	if err != nil {
		return false, err
	}
	flag, exists := values[flagName]
	if exists {
		return flag == true, nil
	}
	return false, domain.ErrKeyNotFound
}

func (s *KoanfStore) loadTargetPath(repoConfig *koanf.Koanf, relativePath string) (domain.Path, error) {
	values, err := s.loadValuesForTemplate(repoConfig, relativePath)
	if err != nil {
		return "", err
	}
	targetPath, exists := values["targetPath"]
	if exists {
		newPath, isString := targetPath.(string)
		if isString {
			if strings.HasSuffix(newPath, "/") {
				return domain.Path(path.Clean(path.Join(newPath, path.Base(relativePath)))), nil
			}
			return domain.Path(newPath), nil
		}
		return "", nil
	}
	return "", nil
}
