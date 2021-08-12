package rendering

import (
	"errors"
	"fmt"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/printer"
)

const (
	// RenderTemplatesEvent identifies an event that will render templates in a local Git repository.
	RenderTemplatesEvent core.EventName = "core:render-templates"
)

type RenderTemplatesHandler struct {
	repoStore     core.GitRepositoryStore
	valueStore    core.ValueStore
	templateStore core.TemplateStore
}

type pipelineContext struct {
	repo      core.GitRepository
	log       printer.Printer
	values    core.Values
	templates []core.Template
}

func NewRenderTemplatesHandler(repoStore core.GitRepositoryStore, ts core.TemplateStore, vs core.ValueStore) *RenderTemplatesHandler {
	return &RenderTemplatesHandler{
		repoStore:     repoStore,
		templateStore: ts,
		valueStore:    vs,
	}
}

func (s *RenderTemplatesHandler) Handle(source core.EventSource) core.EventResult {
	if source.Url == nil {
		return core.EventResult{Error: fmt.Errorf("no URL defined")}
	}
	repo, err := s.repoStore.FetchGitRepository(source.Url)
	if err != nil {
		return core.ToResult(source, err)
	}
	return core.ToResult(source, s.runPipeline(repo))
}

func (s *RenderTemplatesHandler) runPipeline(repo core.GitRepository) error {
	ctx := &pipelineContext{
		repo: repo,
		log:  printer.New().SetName(repo.GetConfig().URL.GetRepositoryName()),
	}

	logger := printer.PipelineLogger{Logger: ctx.log}
	result := pipeline.NewPipelineWithLogger(logger).WithSteps(
		pipeline.NewStep("load templates", s.toAction(ctx, s.loadTemplates)),
		pipeline.NewStep("render templates", s.toAction(ctx, s.renderTemplates)),
		pipeline.NewStep("delete unwanted files", s.toAction(ctx, s.cleanupFiles)),
	).Run()
	return result.Err
}

func (s *RenderTemplatesHandler) toAction(ctx *pipelineContext, action func(ctx *pipelineContext) error) pipeline.ActionFunc {
	return func() pipeline.Result {
		return pipeline.Result{Err: action(ctx)}
	}
}

func (s *RenderTemplatesHandler) renderTemplates(ctx *pipelineContext) error {
	props := ctx.repo.GetConfig()
	for _, template := range ctx.templates {
		if unmanaged, err := s.valueStore.FetchUnmanagedFlag(template, &props); err != nil && !errors.Is(err, core.ErrKeyNotFound) {
			return err
		} else if unmanaged {
			continue
		}
		if err := s.loadValues(ctx, template); err != nil {
			return err
		}
		return s.renderTemplate(ctx, template)
	}
	return nil
}

func (s *RenderTemplatesHandler) renderTemplate(ctx *pipelineContext, template core.Template) error {
	props := ctx.repo.GetConfig()
	alternativePath, err := s.valueStore.FetchTargetPath(template, &props)
	if err != nil {
		return err
	}
	result, err := template.Render(ctx.values)
	if err != nil {
		return err
	}

	targetPath := template.GetRelativePath()
	if alternativePath != "" {
		targetPath = alternativePath
	}

	return ctx.repo.EnsureFile(targetPath, result, template.GetFileMode())
}

func (s *RenderTemplatesHandler) loadTemplates(ctx *pipelineContext) error {
	templates, err := s.templateStore.FetchTemplates()
	if err != nil {
		return err
	}
	ctx.templates = templates
	return nil
}

func (s *RenderTemplatesHandler) loadValues(ctx *pipelineContext, template core.Template) error {
	props := ctx.repo.GetConfig()
	values, err := s.valueStore.FetchValuesForTemplate(template, &props)
	if err != nil {
		return err
	}
	ctx.values = core.Values{
		"Values":   values,
		"Metadata": props,
	}
	return nil
}

func (s *RenderTemplatesHandler) cleanupFiles(ctx *pipelineContext) error {
	props := ctx.repo.GetConfig()
	files, err := s.valueStore.FetchFilesToDelete(&props)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err := ctx.repo.DeleteFile(file); err != nil {
			return err
		}
	}
	return nil
}
