package templateengine

import (
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/printer"
)

type RenderServiceInstrumentation struct {
	log printer.Printer
}

func NewRenderServiceInstrumentation() *RenderServiceInstrumentation {
	i := &RenderServiceInstrumentation{}
	return i
}

func (r *RenderServiceInstrumentation) NewRenderServiceInstrumentation(repository *domain.GitRepository) domain.RenderServiceInstrumentation {
	newCopy := NewRenderServiceInstrumentation()
	newCopy.log = printer.New().SetName(repository.URL.GetFullName()).SetLevel(printer.DefaultLevel)
	return newCopy
}

func (r *RenderServiceInstrumentation) FetchedTemplatesFromStore(fetchErr error) error {
	if fetchErr != nil {
		r.log.DebugF("Fetched templates")
	}
	return fetchErr
}

func (r *RenderServiceInstrumentation) FetchedValuesForTemplate(fetchErr error, template *domain.Template) error {
	if fetchErr == nil {
		r.log.DebugF("Fetched Values for template '%s'", template.RelativePath)
	}
	return fetchErr
}

func (r *RenderServiceInstrumentation) AttemptingToRenderTemplate(template *domain.Template) {
	r.log.DebugF("Rendering template '%s'...", template.RelativePath)
}

func (r *RenderServiceInstrumentation) WrittenRenderResultToFile(template *domain.Template, targetPath domain.Path, writeErr error) error {
	if writeErr == nil {
		r.log.InfoF("Rendered '%s' to '%s'", template.RelativePath, targetPath)
	}
	return writeErr
}
