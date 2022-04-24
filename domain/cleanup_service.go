package domain

import (
	"context"
	"os"

	pipeline "github.com/ccremer/go-command-pipeline"
)

type CleanupService struct {
	instrumentation CleanupServiceInstrumentation
}

type CleanupPipeline struct {
	Repository    *GitRepository
	ValueStore    ValueStore
	TemplateStore TemplateStore

	files     []Path
	templates []*Template

	instrumentation CleanupServiceInstrumentation
}

func NewCleanupService(
	instrumentation CleanupServiceInstrumentation,
) *CleanupService {
	return &CleanupService{
		instrumentation: instrumentation,
	}
}

func (s *CleanupService) CleanupUnwantedFiles(pipe CleanupPipeline) error {
	pipe.instrumentation = s.instrumentation.WithRepository(pipe.Repository)
	result := pipeline.NewPipeline().WithSteps(
		pipeline.NewStepFromFunc("preflight check", pipe.preFlightCheck),
		pipeline.NewStepFromFunc("load templates", pipe.loadTemplates),
		pipeline.NewStepFromFunc("load files", pipe.loadFiles),
		pipeline.NewStepFromFunc("delete files", pipe.deleteFiles),
	).Run()
	return result.Err()
}

func (p *CleanupPipeline) preFlightCheck(_ context.Context) error {
	err := firstOf(
		checkIfArgumentNil(p.Repository, "Repository"),
		checkIfArgumentNil(p.ValueStore, "ValueStore"),
		checkIfArgumentNil(p.TemplateStore, "TemplateStore"),
	)
	return err
}

func (p *CleanupPipeline) loadTemplates(_ context.Context) error {
	templates, err := p.TemplateStore.FetchTemplates()
	p.templates = templates
	return err
}

func (p *CleanupPipeline) loadFiles(_ context.Context) error {
	files, err := p.ValueStore.FetchFilesToDelete(p.Repository, p.templates)
	p.files = files
	return p.instrumentation.FetchedFilesToDelete(err, files)
}

func (p *CleanupPipeline) deleteFiles(_ context.Context) error {
	for _, file := range p.files {
		absoluteFile := p.Repository.RootDir.Join(file)
		if absoluteFile.FileExists() {
			if err := os.Remove(absoluteFile.String()); hasFailed(err) {
				return err
			}
			p.instrumentation.DeletedFile(absoluteFile)
		}
	}
	return nil
}
