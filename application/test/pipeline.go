package test

import (
	"context"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/domain"
)

type updatePipeline struct {
	pipeline.Pipeline
	repo       *domain.GitRepository
	appService *AppService
}

func (c *updatePipeline) renderTemplates(_ context.Context) error {
	err := c.appService.renderService.RenderTemplates(domain.RenderContext{
		Repository:    c.repo,
		ValueStore:    c.appService.valueStore,
		TemplateStore: c.appService.templateStore,
		Engine:        c.appService.engine,
	})
	return err
}

func (c *updatePipeline) diff(_ context.Context) error {
	diff, err := c.appService.repoStore.Diff(c.repo, domain.DiffOptions{
		WorkDirToHEAD: c.appService.cfg.Git.SkipCommit, // If we don't commit, show the unstaged changes
	})
	if err != nil {
		return err
	}
	c.appService.diffPrinter.PrintDiff("Diff: "+c.repo.URL.GetFullName(), diff)
	return nil
}
