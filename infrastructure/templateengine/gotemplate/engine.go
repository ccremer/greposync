package gotemplate

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/ccremer/greposync/domain"
)

type GoTemplateEngine struct {
	RootDir domain.Path

	cache map[domain.Path]*template.Template
}

const (
	// ErrorOnMissingKey is the option that configures template to exit on missing values.
	ErrorOnMissingKey = "missingkey=error"
	// HelperFileName is the base file name that is not considered a template, but provides additional template definitions.
	HelperFileName = "_helpers.tpl"
)

var templateFunctions = GoTemplateFuncMap()

func NewEngine() *GoTemplateEngine {
	return &GoTemplateEngine{
		cache: map[domain.Path]*template.Template{},
	}
}

func (e *GoTemplateEngine) Execute(template *domain.Template, values domain.Values) (domain.RenderResult, error) {
	tpl, err := e.loadGoTemplate(template)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, values)
	return domain.RenderResult(buf.String()), err
}

func (e *GoTemplateEngine) ExecuteString(templateString string, values domain.Values) (domain.RenderResult, error) {
	buf := &bytes.Buffer{}
	tpl, err := template.
		New("").
		Option(ErrorOnMissingKey).
		Funcs(templateFunctions).
		Parse(templateString)
	if err != nil {
		return "", err
	}
	err = tpl.Execute(buf, values)
	return domain.RenderResult(buf.String()), err
}

func (e *GoTemplateEngine) loadGoTemplate(template *domain.Template) (*template.Template, error) {
	if tpl, exists := e.cache[template.RelativePath]; exists {
		return tpl, nil
	}
	fullFilePath := filepath.Join(e.RootDir.String(), template.RelativePath.String())
	helperPath := domain.NewFilePath(e.RootDir.String(), HelperFileName)
	tpl, err := e.parseTemplateFile(fullFilePath, helperPath)
	if err != nil {
		return nil, err
	}
	e.cache[template.RelativePath] = tpl
	return tpl, nil
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
