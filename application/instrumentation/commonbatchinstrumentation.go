package instrumentation

import (
	"fmt"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/go-command-pipeline/parallel"
	"github.com/ccremer/greposync/application/clierror"
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/ccremer/greposync/infrastructure/ui"
	"github.com/go-logr/logr"
	"github.com/hashicorp/go-multierror"
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

func (i *CommonBatchInstrumentation) BatchPipelineStarted(repos []*domain.GitRepository) {
	i.log.Info("Update started")
	p, _ := i.console.BatchProgressbar.WithTotal(len(repos)).Start()
	i.console.BatchProgressbar = p
}

func (i *CommonBatchInstrumentation) BatchPipelineCompleted(repos []*domain.GitRepository) {
	i.log.Info("Update finished")

	for index, result := range i.results {
		if result.IsFailed() {
			repo := repos[index]
			i.console.Flush(repo.URL.GetFullName(), "Log: "+repo.URL.GetFullName())
		}
	}
}

func (i *CommonBatchInstrumentation) PipelineForRepositoryStarted(repo *domain.GitRepository) {
	i.log.WithName(repo.URL.GetFullName()).V(logging.LevelDebug).Info("Starting pipeline")
}

func (i *CommonBatchInstrumentation) PipelineForRepositoryCompleted(repo *domain.GitRepository, err error) {
	if err != nil {
		i.log.WithName(repo.URL.GetFullName()).Error(nil, err.Error())
	}
	i.console.PrintProgressbarMessage(repo.URL.GetFullName(), err)
}

func (i *CommonBatchInstrumentation) NewCollectErrorHandler(repos []*domain.GitRepository, skipBroken bool) parallel.ResultHandler {
	if skipBroken {
		return i.ignoreErrors()
	}
	return i.reduceErrors(repos)
}

func (i *CommonBatchInstrumentation) ignoreErrors() parallel.ResultHandler {
	// Do not propagate update errors from single repositories up the stack
	return func(ctx pipeline.Context, results map[uint64]pipeline.Result) pipeline.Result {
		i.results = results
		return pipeline.Result{}
	}
}

func (i *CommonBatchInstrumentation) reduceErrors(repos []*domain.GitRepository) parallel.ResultHandler {
	return func(ctx pipeline.Context, results map[uint64]pipeline.Result) pipeline.Result {
		i.results = results
		var err error
		for index, repo := range repos {
			if result := results[uint64(index)]; result.Err != nil {
				err = multierror.Append(err, fmt.Errorf("%s: %w", repo.URL.GetRepositoryName(), result.Err))
			}
		}
		return pipeline.Result{Err: fmt.Errorf("%w: %s", clierror.ErrPipeline, err)}
	}
}
