package update

import (
	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/ccremer/greposync/infrastructure/ui"
	"github.com/go-logr/logr"
)

type UpdateInstrumentation struct {
	console *ui.ColoredConsole
	log     logr.Logger

	results map[uint64]pipeline.Result
}

func NewUpdateInstrumentation(console *ui.ColoredConsole, factory logging.LoggerFactory) *UpdateInstrumentation {
	return &UpdateInstrumentation{
		console: console,
		log:     factory.NewGenericLogger(""),
	}
}

func (i *UpdateInstrumentation) batchPipelineStarted(total int) {
	i.log.Info("Update started")
	p, _ := i.console.BatchProgressbar.WithTotal(total).Start()
	i.console.BatchProgressbar = p
}

func (i *UpdateInstrumentation) batchPipelineCompleted(repos []*domain.GitRepository) {
	i.log.Info("Update finished")

	for index, result := range i.results {
		if result.IsFailed() {
			repo := repos[index]
			i.console.Flush(repo.URL.GetFullName(), "Log: "+repo.URL.GetFullName())
		}
	}
}

func (i *UpdateInstrumentation) pipelineForRepositoryStarted(repo *domain.GitRepository) {
	i.log.WithName(repo.URL.GetFullName()).V(logging.LevelDebug).Info("Starting pipeline")
}

func (i *UpdateInstrumentation) pipelineForRepositoryCompleted(repo *domain.GitRepository, err error) {
	if err != nil {
		i.log.WithName(repo.URL.GetFullName()).Error(nil, err.Error())
	}
	i.console.PrintProgressbarMessage(repo.URL.GetFullName(), err)
}
