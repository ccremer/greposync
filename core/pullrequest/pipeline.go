package pullrequest

import (
	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/printer"
)

type pipelineContext struct {
	template core.Template
	repo     core.GitRepository
	body     string
	pr       core.PullRequest
}

func (s *PullRequestService) RunPipeline(repo core.GitRepository) error {
	ctx := &pipelineContext{
		repo: repo,
	}
	s.log.SetName(repo.GetConfig().URL.GetRepositoryName())
	logger := printer.PipelineLogger{Logger: s.log}
	result := pipeline.NewPipelineWithLogger(logger).WithSteps(
		pipeline.NewStep("fetch existing pr", s.toAction(ctx, s.fetchExistingPr)),
		pipeline.NewStep("fetch pr template", s.toAction(ctx, s.fetchPrTemplate)),
		pipeline.NewStep("render template", s.toAction(ctx, s.renderTemplate)),
		pipeline.NewStep("create or update pr", s.toAction(ctx, s.createOrUpdatePr)),
	).Run()
	return result.Err
}

func (s *PullRequestService) toAction(ctx *pipelineContext, action func(ctx *pipelineContext) error) pipeline.ActionFunc {
	return func() pipeline.Result {
		return pipeline.Result{Err: action(ctx)}
	}
}
