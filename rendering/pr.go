package rendering

import (
	"bytes"
	"path"
	"text/template"

	pipeline "github.com/ccremer/go-command-pipeline"
)

// RenderPrTemplate renders the PR template.
// If the BodyTemplate config is a path to an existing file, it will use the file and overwrite the config with the rendered result.
// IF the BodyTemplate config is a string, it will overwrite it with a rendered and data-injected string.
// If the BodyTemplate config is empty, it will use the CommitMessage.
func (r *Renderer) RenderPrTemplate() pipeline.ActionFunc {
	return func() pipeline.Result {
		t := r.cfg.PullRequest.BodyTemplate
		if t == "" {
			r.p.InfoF("No PullRequest template defined")
			t = r.cfg.Git.CommitMessage
		}

		data := Values{"Metadata": r.ConstructMetadata()}
		filePath := path.Clean(t)
		if fileExists(filePath) {
			if str, err := r.renderTemplateFile(data, t); err != nil {
				return pipeline.Result{Err: err}
			} else {
				r.cfg.PullRequest.BodyTemplate = str
			}
		} else {
			if str, err := r.renderString(data, t); err != nil {
				return pipeline.Result{Err: err}
			} else {
				r.cfg.PullRequest.BodyTemplate = str
			}
		}
		return pipeline.Result{}
	}
}

func (r *Renderer) renderString(data interface{}, content string) (string, error) {
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

func (r *Renderer) renderTemplateFile(data interface{}, filePath string) (string, error) {
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
