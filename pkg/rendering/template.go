package rendering

import (
	"bufio"
	"errors"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
	"golang.org/x/sys/unix"
)

type (
	// GoTemplateService implements core.TemplateFacade.
	GoTemplateService struct {
		config    *cfg.TemplateConfig
		templates map[string]*template.Template
	}
)

const (
	ErrorOnMissingKey = "missingkey=error"
	HelperFileName    = "_helpers.tpl"
)

var (
	templateFunctions = GoTemplateFuncMap()
)

func NewTemplateInstance(config *cfg.TemplateConfig) *GoTemplateService {
	return &GoTemplateService{
		config: config,
	}
}

// FetchTemplates implements core.TemplateFacade.
func (s *GoTemplateService) FetchTemplates() ([]core.Template, error) {
	s.templates = make(map[string]*template.Template)
	templates, err := s.listAllTemplates()
	if err != nil {
		return []core.Template{}, err
	}
	for _, file := range templates {
		if err := s.parseTemplate(file); err != nil {
			return []core.Template{}, err
		}
	}
	return templates, nil
}

func (s *GoTemplateService) RenderTemplate(output core.Output) error {
	t, exists := s.templates[output.Template.RelativePath]
	if !exists {
		return errors.New("template does not exist: " + output.Template.RelativePath)
	}
	if err := s.createDirs(path.Join(output.Git.RootDir, output.TargetPath)); err != nil {
		return err
	}
	return s.writeTemplate(output, t)
}

func (s *GoTemplateService) listAllTemplates() (templates []core.Template, err error) {
	err = filepath.Walk(path.Clean(s.config.RootDir),
		func(file string, info os.FileInfo, err error) error {
			tpl, pathErr := s.evaluatePath(file, info, err)
			if pathErr != nil || tpl == nil {
				return pathErr
			}
			templates = append(templates, *tpl)
			return nil
		})
	return templates, err
}

func (s *GoTemplateService) evaluatePath(file string, info os.FileInfo, err error) (*core.Template, error) {
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
	return &core.Template{
		RelativePath: relativePath,
		FileMode:     info.Mode(),
	}, nil
}

func (s *GoTemplateService) parseTemplate(tpl core.Template) error {
	originalFileName := path.Base(tpl.RelativePath)
	originalTemplatePath := path.Join(s.config.RootDir, tpl.RelativePath)

	templates := []string{originalTemplatePath}
	helperFilePath := path.Join(s.config.RootDir, HelperFileName)
	if fileExists(helperFilePath) {
		templates = append(templates, helperFilePath)
	}
	// Read template and helpers
	t, err := template.
		New(originalFileName).
		Option(ErrorOnMissingKey).
		Funcs(templateFunctions).
		ParseFiles(templates...)
	s.templates[tpl.RelativePath] = t
	return err
}

func (s *GoTemplateService) createDirs(targetPath string) error {
	dir := path.Dir(targetPath)
	return os.MkdirAll(dir, 0775)
}

func (s *GoTemplateService) writeTemplate(output core.Output, t *template.Template) error {
	// This allows us to create files with 777 permissions
	originalUmask := unix.Umask(0)
	defer unix.Umask(originalUmask)

	fileName := path.Join(output.Git.RootDir, output.TargetPath)
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, output.Template.FileMode)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	err = t.Execute(w, output.Values)
	if err != nil {
		return err
	}
	return w.Flush()
}

func fileExists(fileName string) bool {
	if info, err := os.Stat(fileName); err == nil && !info.IsDir() {
		return true
	}
	return false
}

// GetTemplateInstances gets the template map created by FetchTemplates.
//
// Deprecated: This method is meant to be temporarily as long as DDD refactoring isn't completed yet.
func (s *GoTemplateService) GetTemplateInstances() map[string]*template.Template {
	return s.templates
}

// SetTemplateInstances sets the internal template map.
//
// Deprecated: This method exists only for testing purposes and will be removed with completed DDD refactoring.
func (s *GoTemplateService) SetTemplateInstances(templates map[string]*template.Template) {
	s.templates = templates
}
