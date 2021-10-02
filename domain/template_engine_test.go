package domain

import (
	"bytes"
	"path"
	"text/template"
)

type DummyEngine struct {
	templatePath string
}

func (d DummyEngine) Execute(templ *Template, values Values) (RenderResult, error) {
	tpl, err := template.New(path.Base(templ.RelativePath.String())).ParseFiles(d.templatePath)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, values)

	return RenderResult(buf.String()), err
}
