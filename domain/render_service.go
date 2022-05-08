package domain

import (
	"context"
	"errors"
	"os"
	"path/filepath"

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

// RenderTemplates loads the TemplateStore and renders them in the GitRepository.RootDir of the given RenderContext.Repository.
func (s *RenderService) RenderTemplates(ctx RenderContext) error {
	ctx.instrumentation = s.instrumentation.WithRepository(ctx.Repository)
	result := pipeline.NewPipeline().WithSteps(
		pipeline.NewStepFromFunc("preflight check", ctx.preFlightCheck),
		pipeline.NewStepFromFunc("load templates", ctx.loadTemplates),
		pipeline.NewStepFromFunc("render templates", ctx.renderTemplates),
	).Run()
	return result.Err()
}

func (ctx *RenderContext) preFlightCheck(_ context.Context) error {
	err := firstOf(
		checkIfArgumentNil(ctx.Engine, "Engine"),
		checkIfArgumentNil(ctx.Repository, "Repository"),
		checkIfArgumentNil(ctx.TemplateStore, "TemplateStore"),
		checkIfArgumentNil(ctx.ValueStore, "ValueStore"),
	)
	return err
}

func (ctx *RenderContext) renderTemplates(_ context.Context) error {
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

	actualFile := ctx.Repository.RootDir.Join(targetPath)
	err = os.MkdirAll(filepath.Dir(actualFile.String()), 0775)
	if err != nil {
		return err
	}

	err = result.WriteToFile(actualFile, template.FilePermissions)
	return ctx.instrumentation.WrittenRenderResultToFile(template, targetPath, err)
}

func (ctx *RenderContext) loadTemplates(_ context.Context) error {
	templates, err := ctx.TemplateStore.FetchTemplates()
	ctx.templates = templates
	return ctx.instrumentation.FetchedTemplatesFromStore(err)
}

func (ctx *RenderContext) loadValues(template *Template) error {
	values, err := ctx.ValueStore.FetchValuesForTemplate(template, ctx.Repository)
	ctx.values = ctx.enrichWithMetadata(values, template)
	return ctx.instrumentation.FetchedValuesForTemplate(err, template)
}

func (ctx *RenderContext) enrichWithMetadata(values Values, template *Template) Values {
	return Values{
		ValuesKey: values,
		MetadataValueKey: Values{
			RepositoryValueKey: ctx.Repository.AsValues(),
			TemplateValueKey:   template.AsValues(),
		},
	}
}
