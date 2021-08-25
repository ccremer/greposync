package update

import (
	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/go-command-pipeline/predicate"
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/printer"
)

type pipelineContext struct {
	log           printer.Printer
	repo          *domain.GitRepository
	renderService *domain.RenderService
	appService    *AppService
	differ        *Differ
}

func (c *pipelineContext) clone() pipeline.ActionFunc {
	return c.toAction(c.appService.repoStore.Clone)
}

func (c *pipelineContext) fetch() pipeline.ActionFunc {
	return c.toAction(c.appService.repoStore.Fetch)
}

func (c *pipelineContext) pull() pipeline.ActionFunc {
	return c.toAction(c.appService.repoStore.Pull)
}

func (c *pipelineContext) checkout() pipeline.ActionFunc {
	return c.toAction(c.appService.repoStore.Checkout)
}

func (c *pipelineContext) reset() pipeline.ActionFunc {
	return c.toAction(c.appService.repoStore.Reset)
}

func (c *pipelineContext) add() pipeline.ActionFunc {
	return c.toAction(c.appService.repoStore.Add)
}

func (c *pipelineContext) commit() pipeline.ActionFunc {
	return func() pipeline.Result {
		err := c.appService.repoStore.Commit(c.repo, domain.CommitOptions{
			Message: "asdf",
			Amend:   true,
		})
		return pipeline.Result{Err: err}
	}
}

func (c *pipelineContext) diff() pipeline.ActionFunc {
	return func() pipeline.Result {
		diff, err := c.appService.repoStore.Diff(c.repo)
		if err != nil {
			return pipeline.Result{Err: err}
		}
		c.differ.PrettyPrint(diff)
		return pipeline.Result{}
	}
}

func (c *pipelineContext) renderTemplates() pipeline.ActionFunc {
	return func() pipeline.Result {
		err := c.renderService.RenderTemplates(domain.RenderContext{
			Repository:    c.repo,
			ValueStore:    c.appService.valueStore,
			TemplateStore: c.appService.templateStore,
			Engine:        c.appService.engine,
		})
		return pipeline.Result{Err: err}
	}
}

func (c *pipelineContext) dirMissing() predicate.Predicate {
	return func(step pipeline.Step) bool {
		return !c.repo.RootDir.DirExists()
	}
}

func (c *pipelineContext) toAction(f func(repository *domain.GitRepository) error) pipeline.ActionFunc {
	return func() pipeline.Result {
		err := f(c.repo)
		return pipeline.Result{Err: err}
	}
}
