package rendering

import (
	"bytes"
	"path"
	"text/template"
)

func (s *Service) RenderPrTemplate(data Data) string {
	fileName := path.Base(s.cfg.PrTemplatePath)
	tpl, err := template.
		New(fileName).
		Option(errorOnMissingKey).
		Funcs(templateFunctions).
		ParseFiles(s.cfg.PrTemplatePath)
	s.p.CheckIfError(err)
	buf := bytes.NewBuffer([]byte{})
	err = tpl.Execute(buf, data)
	s.p.CheckIfError(err)
	return buf.String()
}
