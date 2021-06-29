package rendering

import (
	"bytes"
	"path"
	"text/template"
)

func (r *Renderer) RenderString(data interface{}, content string) string {
	r.p.DebugF("Parsing template from string")
	tpl, err := template.
		New("").
		Option(errorOnMissingKey).
		Funcs(templateFunctions).
		Parse(content)
	r.p.CheckIfError(err)
	buf := bytes.NewBuffer([]byte{})
	err = tpl.Execute(buf, data)
	r.p.CheckIfError(err)
	return buf.String()
}

func (r *Renderer) RenderTemplateFile(data interface{}, filePath string) string {
	r.p.DebugF("Parsing template file %s", filePath)
	fileName := path.Base(filePath)
	tpl, err := template.
		New(fileName).
		Option(errorOnMissingKey).
		Funcs(templateFunctions).
		ParseFiles(filePath)
	r.p.CheckIfError(err)
	buf := bytes.NewBuffer([]byte{})
	err = tpl.Execute(buf, data)
	r.p.CheckIfError(err)
	return buf.String()
}
