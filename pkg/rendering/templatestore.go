package rendering

import (
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
)

// GoTemplateStore implements core.TemplateStore.
type GoTemplateStore struct {
	config *cfg.TemplateConfig
}

const (
	// ErrorOnMissingKey is the option that configures template to exit on missing values.
	ErrorOnMissingKey = "missingkey=error"
	// HelperFileName is the base file name that is not considered a template, but provides additional template definitions.
	HelperFileName = "_helpers.tpl"
)

var (
	templateFunctions = GoTemplateFuncMap()
)

// NewGoTemplateStore returns a new GoTemplateStore instance.
func NewGoTemplateStore(config *cfg.TemplateConfig) *GoTemplateStore {
	return &GoTemplateStore{
		config: config,
	}
}

// FetchTemplates implements core.TemplateStore.
func (s *GoTemplateStore) FetchTemplates() ([]core.Template, error) {
	templates, err := s.listAllTemplates()
	if err != nil {
		return []core.Template{}, err
	}
	for _, tpl := range templates {
		t, err := s.parseTemplate(tpl)
		if err != nil {
			return []core.Template{}, err
		}
		tpl.template = t
		templates = append(templates, tpl)
	}

	// convert to interface type
	list := make([]core.Template, len(templates))
	for index := range templates {
		list[index] = templates[index]
	}
	return list, nil
}

func (s *GoTemplateStore) listAllTemplates() (templates []*GoTemplate, err error) {
	err = filepath.Walk(path.Clean(s.config.RootDir),
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

func (s *GoTemplateStore) evaluatePath(file string, info os.FileInfo, err error) (*GoTemplate, error) {
	if err != nil {
		return nil, err
	}
	// Don't add helper file or directories
	if path.Base(file) == HelperFileName || info.IsDir() {
		return nil, nil
	}
	relativePath, pathErr := filepath.Rel(s.config.RootDir, file)
	if pathErr != nil {
		return nil, pathErr
	}
	return &GoTemplate{
		RelativePath: relativePath,
		FileMode:     info.Mode(),
	}, nil
}

func (s *GoTemplateStore) parseTemplate(tpl *GoTemplate) (*template.Template, error) {
	originalFileName := path.Base(tpl.RelativePath)
	originalTemplatePath := path.Join(s.config.RootDir, tpl.RelativePath)

	templates := []string{originalTemplatePath}
	helperFilePath := path.Join(s.config.RootDir, HelperFileName)
	if fileExists(helperFilePath) {
		templates = append(templates, helperFilePath)
	}
	// Read template and helpers
	return template.
		New(originalFileName).
		Option(ErrorOnMissingKey).
		Funcs(templateFunctions).
		ParseFiles(templates...)
}

func fileExists(fileName string) bool {
	if info, err := os.Stat(fileName); err == nil && !info.IsDir() {
		return true
	}
	return false
}
