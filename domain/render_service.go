package domain

import (
	"errors"
	"os"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/core"
)

type RenderService struct{}

type RenderContext struct {
	Repository    *GitRepository
	ValueStore    ValueStore
	TemplateStore TemplateStore
	Engine        TemplateEngine

	templates []*Template
	values    Values
}

func NewRenderService() *RenderService {
	return &RenderService{}
}

func (s *RenderService) RenderTemplates(ctx RenderContext) error {
	result := pipeline.NewPipeline().WithSteps(
		pipeline.NewStep("preflight check", ctx.preFlightCheck()),
		pipeline.NewStep("load templates", ctx.toAction(ctx.loadTemplates)),
		pipeline.NewStep("render templates", ctx.toAction(ctx.renderTemplates)),
	).Run()
	return result.Err
}

func (ctx *RenderContext) preFlightCheck() pipeline.ActionFunc {
	return func() pipeline.Result {
		err := firstOf(
			checkIfArgumentNil(ctx.Engine, "Engine"),
			checkIfArgumentNil(ctx.Repository, "Repository"),
			checkIfArgumentNil(ctx.TemplateStore, "TemplateStore"),
			checkIfArgumentNil(ctx.ValueStore, "ValueStore"),
		)
		return pipeline.Result{Err: err}
	}
}

func (ctx *RenderContext) toAction(action func() error) pipeline.ActionFunc {
	return func() pipeline.Result {
		return pipeline.Result{Err: action()}
	}
}

func (ctx *RenderContext) renderTemplates() error {
	for _, template := range ctx.templates {
		if unmanaged, err := ctx.ValueStore.FetchUnmanagedFlag(template, ctx.Repository); err != nil && !errors.Is(err, core.ErrKeyNotFound) {
			return err
		} else if unmanaged {
			continue
		}
		if err := ctx.loadValues(template); err != nil {
			return err
		}
		return ctx.renderTemplate(template)
	}
	return nil
}

func (ctx *RenderContext) renderTemplate(template *Template) error {
	alternativePath, err := ctx.ValueStore.FetchTargetPath(template, ctx.Repository)
	if err != nil {
		return err
	}
	result, err := template.Render(ctx.values, ctx.Engine)
	if err != nil {
		return err
	}

	targetPath := template.RelativePath
	if alternativePath != "" {
		targetPath = alternativePath
	}

	return os.WriteFile(ctx.Repository.RootDir.Join(targetPath).String(), []byte(result), template.FilePermissions.FileMode())
}

func (ctx *RenderContext) loadTemplates() error {
	templates, err := ctx.TemplateStore.FetchTemplates()
	ctx.templates = templates
	return err
}

func (ctx *RenderContext) loadValues(template *Template) error {
	values, err := ctx.ValueStore.FetchValuesForTemplate(template, ctx.Repository)
	ctx.values = Values{
		"Values":   values,
		"Metadata": ctx.Repository,
	}
	return err
}
