package gitrepo

import (
	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/go-command-pipeline/predicate"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/printer"
)

const (
	// PrepareWorkspaceEvent identifies an event that will prepare a local Git repository.
	PrepareWorkspaceEvent core.EventName = "core:prepare-workspace"
)

type pipelineContext struct {
	repo core.GitRepository
}

func (s *PrepareWorkspaceHandler) runPipeline(repo core.GitRepository) error {
	ctx := &pipelineContext{
		repo: repo,
	}

	gitDirExists := s.dirExists(repo.GetConfig().RootDir)
	// TODO: make reset configurable
	resetEnabled := true
	s.log.SetName(repo.GetConfig().URL.GetRepositoryName())
	logger := printer.PipelineLogger{Logger: s.log}
	result := pipeline.NewPipelineWithLogger(logger).WithSteps(
		predicate.ToStep("clone repository", s.toAction(ctx, s.clone), If(!gitDirExists)),
		predicate.ToStep("fetch", s.toAction(ctx, s.fetch), If(resetEnabled)),
		predicate.ToStep("reset repository", s.toAction(ctx, s.reset), If(resetEnabled)),
		pipeline.NewStep("checkout branch", s.toAction(ctx, s.checkout)),
		predicate.ToStep("pull", s.toAction(ctx, s.pull), If(resetEnabled)),
	).Run()
	return result.Err
}

func (s *PrepareWorkspaceHandler) toAction(ctx *pipelineContext, action func(ctx *pipelineContext) error) pipeline.ActionFunc {
	return func() pipeline.Result {
		return pipeline.Result{Err: action(ctx)}
	}
}

func If(v bool) predicate.Predicate {
	return func(step pipeline.Step) bool {
		return v
	}
}
