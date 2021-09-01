package gotemplate

import (
	"os"
	"path/filepath"

	"github.com/ccremer/greposync/domain"
)

type GoTemplateStore struct {
	RootDir string
}

func NewTemplateStore() *GoTemplateStore {
	return &GoTemplateStore{}
}

func (s *GoTemplateStore) FetchTemplates() ([]*domain.Template, error) {
	templates, err := s.listAllTemplates()
	return templates, err
}

func (s *GoTemplateStore) listAllTemplates() (templates []*domain.Template, err error) {
	err = filepath.Walk(filepath.Clean(s.RootDir),
		func(file string, info os.FileInfo, err error) error {
			tpl, pathErr := s.evaluatePath(file, info, err)
			if pathErr != nil || tpl == nil {
				return pathErr
			}
			templates = append(templates, tpl)
			return nil
		})
	return templates, err
}

func (s *GoTemplateStore) evaluatePath(file string, info os.FileInfo, err error) (*domain.Template, error) {
	if err != nil {
		return nil, err
	}
	// Don't add helper file or directories
	if filepath.Base(file) == HelperFileName || info.IsDir() {
		return nil, nil
	}
	relativePath, pathErr := filepath.Rel(s.RootDir, file)
	if pathErr != nil {
		return nil, pathErr
	}
	return &domain.Template{
		RelativePath:    domain.NewPath(relativePath),
		FilePermissions: domain.Permissions(info.Mode()),
	}, nil
}
