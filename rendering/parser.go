package rendering

import (
	"os"
	"path"
	"path/filepath"
	"text/template"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/cfg"
)

type (
	Parser struct {
		cfg       *cfg.TemplateConfig
		templates map[string]*template.Template
	}
)

const (
	errorOnMissingKey = "missingkey=error"
	HelperFileName    = "_helpers.tpl"
)

// NewParser returns a new reusable parser instance.
func NewParser(cfg *cfg.TemplateConfig) *Parser {
	return &Parser{
		cfg:       cfg,
		templates: map[string]*template.Template{},
	}
}

// ParseTemplateDirAction encapsulates ParseTemplateDir in a pipeline action.
func (r *Parser) ParseTemplateDirAction() pipeline.ActionFunc {
	return func() pipeline.Result {
		return pipeline.Result{Err: r.ParseTemplateDir()}
	}
}

// ParseTemplateDir searches for template files in the template directory and reads them to memory for later execution.
func (r *Parser) ParseTemplateDir() error {
	files, err := r.listAllFiles(path.Clean(r.cfg.RootDir))
	if err != nil {
		return err
	}
	for _, file := range files {
		if err = r.parseTemplate(file); err != nil {
			return err
		}
	}
	return nil
}

func (r *Parser) listAllFiles(root string) (files []string, err error) {
	err = filepath.Walk(root,
		func(file string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fileName := path.Base(file)
			// Don't add helper file or directories
			if fileName != HelperFileName && fileExists(file) {
				files = append(files, file)
			}
			return nil
		})
	return files, err
}

func (r *Parser) parseTemplate(originalTemplatePath string) error {
	originalFileName := path.Base(originalTemplatePath)

	templates := []string{originalTemplatePath}
	helperFilePath := path.Join(r.cfg.RootDir, HelperFileName)
	if fileExists(helperFilePath) {
		templates = append(templates, helperFilePath)
	}
	// Read template and helpers
	tpl, err := template.
		New(originalFileName).
		Option(errorOnMissingKey).
		Funcs(templateFunctions).
		ParseFiles(templates...)
	r.templates[originalTemplatePath] = tpl
	return err
}
