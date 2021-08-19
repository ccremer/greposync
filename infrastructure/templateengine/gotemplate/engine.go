package gotemplate

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/ccremer/greposync/domain"
)

type GoTemplateEngine struct {
	RootDir domain.Path
}

const (
	// ErrorOnMissingKey is the option that configures template to exit on missing values.
	ErrorOnMissingKey = "missingkey=error"
	// HelperFileName is the base file name that is not considered a template, but provides additional template definitions.
	HelperFileName = "_helpers.tpl"
)

var templateFunctions = GoTemplateFuncMap()

func NewEngine() *GoTemplateEngine {
	return &GoTemplateEngine{}
}

func (e *GoTemplateEngine) Execute(template *domain.Template, values domain.Values) (string, error) {
	fileName := filepath.Join(e.RootDir.String(), template.RelativePath.String())
	helperPath := domain.NewFilePath(e.RootDir.String(), HelperFileName)
	tpl, err := e.parseTemplateFile(fileName, helperPath)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, values)
	return buf.String(), err
}

func (e *GoTemplateEngine) parseTemplateFile(fileName string, helperFilePath domain.Path) (*template.Template, error) {
	originalFileName := filepath.Base(fileName)

	templates := []string{fileName}
	if helperFilePath.FileExists() {
		templates = append(templates, helperFilePath.String())
	}
	// Read template and helpers
	return template.
		New(originalFileName).
		Option(ErrorOnMissingKey).
		Funcs(templateFunctions).
		ParseFiles(templates...)
}
