package instrumentation

import (
	"context"
	"errors"
	"fmt"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/application/clierror"
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/ccremer/greposync/infrastructure/ui"
	"github.com/go-logr/logr"
	"github.com/hashicorp/go-multierror"
	"github.com/urfave/cli/v2"
)

type CommonBatchInstrumentation struct {
	console *ui.ColoredConsole
	log     logr.Logger

	results map[uint64]pipeline.Result
}

func NewUpdateInstrumentation(console *ui.ColoredConsole, factory logging.LoggerFactory) *CommonBatchInstrumentation {
	return &CommonBatchInstrumentation{
		console: console,
		log:     factory.NewGenericLogger(""),
	}
}

func (i *CommonBatchInstrumentation) BatchPipelineStarted(message string, repos []*domain.GitRepository) {
	i.log.Info(message)
	i.console.StartBatchUpdate(repos)
}

func (i *CommonBatchInstrumentation) BatchPipelineCompleted(message string, repos []*domain.GitRepository) {
	i.log.Info(message)

	for index, result := range i.results {
		if result.IsFailed() {
			repo := repos[index]
			i.console.Flush(repo.URL.GetFullName(), "Log: "+repo.URL.GetFullName())
		}
	}
}

func (i *CommonBatchInstrumentation) PipelineForRepositoryStarted(repo *domain.GitRepository) {
	i.log.WithName(repo.URL.GetFullName()).V(1).Info("Starting pipeline")
}

func (i *CommonBatchInstrumentation) PipelineForRepositoryCompleted(repo *domain.GitRepository, err error) {
	if err != nil {
		i.log.WithName(repo.URL.GetFullName()).Error(nil, err.Error())
	}
	i.console.PrintProgressbarMessage(repo.URL.GetFullName(), err)
}

func (i *CommonBatchInstrumentation) NewCollectErrorHandler(skipBroken bool) pipeline.ParallelResultHandler {
	if skipBroken {
		return i.ignoreErrors()
	}
	return i.reduceErrors()
}

func (i *CommonBatchInstrumentation) ignoreErrors() pipeline.ParallelResultHandler {
	// Do not propagate update errors from single repositories up the stack, unless one explicitly contains a wrapped exit code error.
	return func(ctx context.Context, results map[uint64]pipeline.Result) error {
		i.results = results
		for _, result := range results {
			if result.IsFailed() {
				var clierr cli.ExitCoder
				if errors.As(result.Err(), &clierr) {
					return clierr
				}
			}
		}
		return nil
	}
}

func (i *CommonBatchInstrumentation) reduceErrors() pipeline.ParallelResultHandler {
	return func(ctx context.Context, results map[uint64]pipeline.Result) error {
		i.results = results
		var err error
		if repos, found := pipeline.LoadFromContext(ctx, RepositoriesContextKey{}); found {
			for index, repo := range repos.([]*domain.GitRepository) {
				if result := results[uint64(index)]; result.Err() != nil {
					err = multierror.Append(err, fmt.Errorf("%s: %w", repo.URL.GetRepositoryName(), result.Err()))
				}
			}
		}
		if err != nil {
			return fmt.Errorf("%w: %s", clierror.ErrPipeline, err)
		}
		return nil
	}
}
