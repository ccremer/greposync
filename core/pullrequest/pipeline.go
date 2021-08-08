package pullrequest

import (
	"fmt"

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

// EnsurePullRequestEvent identifies the event that will create or update a pull request on a remote repository.
const EnsurePullRequestEvent core.EventName = "core:ensure-pullrequest"

func (s *PullRequestHandler) Handle(source core.EventSource) core.EventResult {
	if source.Url == nil {
		return core.EventResult{Error: fmt.Errorf("no URL defined")}
	}
	repo, err := s.repoStore.FetchGitRepository(source.Url)
	if err != nil {
		return core.ToResult(source, err)
	}
	return core.ToResult(source, s.runPipeline(repo))
}

func (s *PullRequestHandler) runPipeline(repo core.GitRepository) error {
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

func (s *PullRequestHandler) toAction(ctx *pipelineContext, action func(ctx *pipelineContext) error) pipeline.ActionFunc {
	return func() pipeline.Result {
		return pipeline.Result{Err: action(ctx)}
	}
}
