package rendering

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
)

// GoTemplateStore implements core.TemplateStore.
type GoTemplateStore struct {
	config *cfg.Configuration
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
func NewGoTemplateStore(config *cfg.Configuration) *GoTemplateStore {
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
	helperPath := filepath.Join(s.config.Template.RootDir, HelperFileName)
	for _, tpl := range templates {
		t, err := s.parseTemplateFile(
			filepath.Join(s.config.Template.RootDir, tpl.RelativePath),
			helperPath)
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

// FetchPullRequestTemplate implements core.TemplateStore.
func (s *GoTemplateStore) FetchPullRequestTemplate() (core.Template, error) {
	t := s.config.PullRequest.BodyTemplate
	if t == "" {
		return nil, nil
	}

	filePath := filepath.Clean(t)
	if fileExists(filePath) {
		tpl, err := s.parseTemplateFile(filePath, "")
		return &GoTemplate{
			template: tpl,
		}, err
	}
	tpl, err := s.parseTemplateString(t)
	return &GoTemplate{
		template: tpl,
	}, err
}

func (s *GoTemplateStore) listAllTemplates() (templates []*GoTemplate, err error) {
	err = filepath.Walk(filepath.Clean(s.config.Template.RootDir),
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
	if filepath.Base(file) == HelperFileName || info.IsDir() {
		return nil, nil
	}
	relativePath, pathErr := filepath.Rel(s.config.Template.RootDir, file)
	if pathErr != nil {
		return nil, pathErr
	}
	return &GoTemplate{
		RelativePath: relativePath,
		FileMode:     info.Mode(),
	}, nil
}

func (s *GoTemplateStore) parseTemplateFile(fileName, helperFilePath string) (*template.Template, error) {
	originalFileName := filepath.Base(fileName)

	templates := []string{fileName}
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

func (s *GoTemplateStore) parseTemplateString(content string) (*template.Template, error) {
	// Read template and helpers
	return template.
		New("").
		Option(ErrorOnMissingKey).
		Funcs(templateFunctions).
		Parse(content)
}

func fileExists(fileName string) bool {
	if info, err := os.Stat(fileName); err == nil && !info.IsDir() {
		return true
	}
	return false
}
