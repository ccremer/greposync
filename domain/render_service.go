package domain

import (
	"errors"

	pipeline "github.com/ccremer/go-command-pipeline"
	"golang.org/x/sys/unix"
)

// RenderService is a domain service that helps rendering templates.
type RenderService struct {
	instrumentation RenderServiceInstrumentation
}

// RenderContext represents a single rendering context for a GitRepository.
type RenderContext struct {
	Repository    *GitRepository
	ValueStore    ValueStore
	TemplateStore TemplateStore
	Engine        TemplateEngine

	instrumentation RenderServiceInstrumentation
	templates       []*Template
	values          Values
}

func NewRenderService(instrumentation RenderServiceInstrumentation) *RenderService {
	return &RenderService{
		instrumentation: instrumentation,
	}
}

// RenderTemplates loads the Templates and renders them in the GitRepository.RootDir of the given RenderContext.Repository.
func (s *RenderService) RenderTemplates(ctx RenderContext) error {
	ctx.instrumentation = s.instrumentation.WithRepository(ctx.Repository)
	result := pipeline.NewPipeline().WithSteps(
		pipeline.NewStep("preflight check", ctx.preFlightCheck()),
		pipeline.NewStep("load templates", ctx.toAction(ctx.loadTemplates)),
		pipeline.NewStep("render templates", ctx.toAction(ctx.renderTemplates)),
	).Run()
	return result.Err
}

func (ctx *RenderContext) preFlightCheck() pipeline.ActionFunc {
	return func(_ pipeline.Context) pipeline.Result {
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
	return func(_ pipeline.Context) pipeline.Result {
		return pipeline.Result{Err: action()}
	}
}

func (ctx *RenderContext) renderTemplates() error {
	for _, template := range ctx.templates {
		if unmanaged, err := ctx.ValueStore.FetchUnmanagedFlag(template, ctx.Repository); err != nil && !errors.Is(err, ErrKeyNotFound) {
			return err
		} else if unmanaged {
			continue
		}
		if err := ctx.loadValues(template); err != nil {
			return err
		}
		if err := ctx.renderTemplate(template); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *RenderContext) renderTemplate(template *Template) error {
	// This allows us to create files with 777 permissions
	originalUmask := unix.Umask(0)
	defer unix.Umask(originalUmask)

	ctx.instrumentation.AttemptingToRenderTemplate(template)
	alternativePath, err := ctx.ValueStore.FetchTargetPath(template, ctx.Repository)
	if err != nil {
		return err
	}
	result, err := template.Render(ctx.values, ctx.Engine)
	if err != nil {
		return err
	}

	targetPath := template.CleanPath()
	if alternativePath != "" {
		targetPath = alternativePath
	}

	err = result.WriteToFile(ctx.Repository.RootDir.Join(targetPath), template.FilePermissions)
	return ctx.instrumentation.WrittenRenderResultToFile(template, targetPath, err)
}

func (ctx *RenderContext) loadTemplates() error {
	templates, err := ctx.TemplateStore.FetchTemplates()
	ctx.templates = templates
	return ctx.instrumentation.FetchedTemplatesFromStore(err)
}

func (ctx *RenderContext) loadValues(template *Template) error {
	values, err := ctx.ValueStore.FetchValuesForTemplate(template, ctx.Repository)
	ctx.values = Values{
		"Values":   values,
		"Metadata": ctx.Repository,
	}
	return ctx.instrumentation.FetchedValuesForTemplate(err, template)
}
