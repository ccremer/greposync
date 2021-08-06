package labels

import (
	"fmt"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/go-command-pipeline/parallel"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/printer"
	"github.com/hashicorp/go-multierror"
)

// RunPipeline implements core.CoreService.
func (s *LabelService) RunPipeline() error {
	logger := printer.PipelineLogger{Logger: s.log}
	result := pipeline.NewPipelineWithLogger(logger).WithSteps(
		pipeline.NewStep("load config", s.loadManagedReposAction()),
		parallel.NewWorkerPoolStep("update labels", 1, s.updateReposInParallel(), s.errorHandler()),
	).Run()
	return result.Err
}

func (s *LabelService) loadManagedReposAction() pipeline.ActionFunc {
	return func() pipeline.Result {
		repos, err := s.repoProvider.FetchGitRepositories()
		s.repoFacades = repos
		return pipeline.Result{Err: err}
	}
}

func (s *LabelService) updateReposInParallel() parallel.PipelineSupplier {
	return func(pipelines chan *pipeline.Pipeline) {
		defer close(pipelines)
		for _, repoFacade := range s.repoFacades {
			p := s.createPipelineForUpdatingLabels(repoFacade)
			pipelines <- p
		}
	}
}

func (s *LabelService) createPipelineForUpdatingLabels(rf core.GitRepository) *pipeline.Pipeline {
	log := printer.New().SetName(rf.GetConfig().URL.GetRepositoryName())
	logger := printer.PipelineLogger{Logger: log}

	p := pipeline.NewPipelineWithLogger(logger)
	p.WithSteps(
		pipeline.NewStep("edit labels", s.createOrUpdateLabelAction(rf)),
		pipeline.NewStep("delete labels", s.deleteLabelAction(rf)),
	)
	return p
}

func (s *LabelService) errorHandler() parallel.ResultHandler {
	return func(results map[uint64]pipeline.Result) pipeline.Result {
		var err error
		for index, service := range s.repoFacades {
			if result := results[uint64(index)]; result.Err != nil {
				err = multierror.Append(err, fmt.Errorf("%s: %w", service.GetConfig().URL.GetRepositoryName(), result.Err))
			}
		}
		return pipeline.Result{Err: err}
	}
}

func (s *LabelService) createOrUpdateLabelAction(r core.GitRepository) pipeline.ActionFunc {
	return func() pipeline.Result {
		return pipeline.Result{Err: s.createOrUpdateLabels(r)}
	}
}

func (s *LabelService) deleteLabelAction(r core.GitRepository) pipeline.ActionFunc {
	return func() pipeline.Result {
		return pipeline.Result{Err: s.deleteLabels(r)}
	}
}
