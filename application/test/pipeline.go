package test

import (
	"context"
	"fmt"
	"os"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/domain"
	"github.com/urfave/cli/v2"
)

type updatePipeline struct {
	pipeline.Pipeline
	repo               *domain.GitRepository
	appService         *AppService
	failPipelineIfDiff bool
}

func (c *updatePipeline) renderTemplates(_ context.Context) error {
	err := c.appService.renderService.RenderTemplates(domain.RenderContext{
		Repository:           c.repo,
		ValueStore:           c.appService.valueStore,
		TemplateStore:        c.appService.templateStore,
		Engine:               c.appService.engine,
		SkipExtensionRemoval: true,
	})
	return err
}

func (c *updatePipeline) diff(_ context.Context) error {
	diff, err := c.appService.repoStore.Diff(c.repo, domain.DiffOptions{})
	if err != nil {
		return err
	}
	c.appService.diffPrinter.PrintDiff("Diff: "+c.repo.URL.GetFullName(), diff)
	if diff != "" && c.failPipelineIfDiff {
		return cli.Exit(fmt.Errorf("diff not empty"), 3)
	}
	return nil
}

func (c *updatePipeline) createOutputDir(_ context.Context) error {
	return os.MkdirAll(c.repo.RootDir.String(), 0775)
}

func (c *updatePipeline) copySyncFile(_ context.Context) error {
	return c.appService.repoStore.CopySyncFile(c.repo)
}
