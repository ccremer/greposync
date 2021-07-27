package labels

import (
	"fmt"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/go-command-pipeline/parallel"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/printer"
	"github.com/hashicorp/go-multierror"
)

func (s *LabelService) RunPipeline() error {
	logger := printer.PipelineLogger{Logger: s.log}
	result := pipeline.NewPipelineWithLogger(logger).WithSteps(
		pipeline.NewStep("load config", s.loadManagedReposAction()),
		pipeline.NewStep("init hosting providers", s.initHostingAPIAction()),
		parallel.NewWorkerPoolStep("update labels", 1, s.updateReposInParallel(), s.errorHandler()),
	).Run()
	return result.Err
}

func (s *LabelService) loadManagedReposAction() pipeline.ActionFunc {
	return func() pipeline.Result {
		repos, err := s.repoProvider.LoadManagedRepositories()
		s.repoFacades = repos
		return pipeline.Result{Err: err}
	}
}

func (s *LabelService) updateReposInParallel() parallel.PipelineSupplier {
	return func(pipelines chan *pipeline.Pipeline) {
		defer close(pipelines)
		for _, repoFacade := range s.repoFacades {
			if hostingProvider, isSupported := s.repoProvider.GetSupportedGitHostingProviders()[repoFacade.GetConfig().Provider]; isSupported {
				p := s.createPipelineForUpdatingLabels(repoFacade, hostingProvider)
				pipelines <- p
			}
		}
	}
}

func (s *LabelService) createPipelineForUpdatingLabels(rf core.GitRepositoryFacade, hf core.GitHostingFacade) *pipeline.Pipeline {
	log := printer.New().SetName(rf.GetConfig().URL.GetRepositoryName())
	logger := printer.PipelineLogger{Logger: log}

	p := pipeline.NewPipelineWithLogger(logger)
	p.WithSteps(
		pipeline.NewStep("edit labels", s.createOrUpdateLabelAction(rf, hf)),
		pipeline.NewStep("delete labels", s.deleteLabelAction(rf, hf)),
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

func (s *LabelService) initHostingAPIAction() pipeline.ActionFunc {
	return func() pipeline.Result {
		return pipeline.Result{Err: s.initHostingAPIs()}
	}
}

func (s *LabelService) createOrUpdateLabelAction(r core.GitRepositoryFacade, h core.GitHostingFacade) pipeline.ActionFunc {
	return func() pipeline.Result {
		return pipeline.Result{Err: s.createOrUpdateLabels(r, h)}
	}
}

func (s *LabelService) deleteLabelAction(r core.GitRepositoryFacade, h core.GitHostingFacade) pipeline.ActionFunc {
	return func() pipeline.Result {
		return pipeline.Result{Err: s.deleteLabels(r, h)}
	}
}
