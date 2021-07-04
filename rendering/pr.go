package rendering

import (
	"bytes"
	"path"
	"text/template"
)

// RenderString renders the template from given content and injected with data.
func (r *Renderer) RenderString(data interface{}, content string) (string, error) {
	r.p.DebugF("Parsing template from string")
	tpl, err := template.
		New("").
		Option(errorOnMissingKey).
		Funcs(templateFunctions).
		Parse(content)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer([]byte{})
	err = tpl.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderTemplateFile renders a template from the given filePath and injected with data.
func (r *Renderer) RenderTemplateFile(data interface{}, filePath string) (string, error) {
	r.p.DebugF("Parsing template file %s", filePath)
	fileName := path.Base(filePath)
	tpl, err := template.
		New(fileName).
		Option(errorOnMissingKey).
		Funcs(templateFunctions).
		ParseFiles(filePath)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer([]byte{})
	err = tpl.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
