package templateengine

import (
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/go-logr/logr"
)

type RenderServiceInstrumentation struct {
	factory logging.LoggerFactory
	log     logr.Logger
}

func NewRenderServiceInstrumentation(factory logging.LoggerFactory) *RenderServiceInstrumentation {
	i := &RenderServiceInstrumentation{
		factory: factory,
	}
	return i
}

func (r *RenderServiceInstrumentation) WithRepository(repository *domain.GitRepository) domain.RenderServiceInstrumentation {
	newCopy := NewRenderServiceInstrumentation(r.factory)
	newCopy.log = r.factory.NewRepositoryLogger(repository)
	return newCopy
}

func (r *RenderServiceInstrumentation) FetchedTemplatesFromStore(fetchErr error) error {
	if fetchErr != nil {
		r.log.V(logging.LevelDebug).Info("Fetched templates")
	}
	return fetchErr
}

func (r *RenderServiceInstrumentation) FetchedValuesForTemplate(fetchErr error, template *domain.Template) error {
	if fetchErr == nil {
		r.log.V(logging.LevelDebug).Info("Fetched Values", "template", template.RelativePath)
	}
	return fetchErr
}

func (r *RenderServiceInstrumentation) AttemptingToRenderTemplate(template *domain.Template) {
	r.log.V(logging.LevelDebug).Info("Rendering template...", "template", template.RelativePath)
}

func (r *RenderServiceInstrumentation) WrittenRenderResultToFile(template *domain.Template, targetPath domain.Path, writeErr error) error {
	if writeErr == nil {
		r.log.V(logging.LevelDebug).Info("Rendered file", "template", template.RelativePath, "target", targetPath)
	}
	return writeErr
}
