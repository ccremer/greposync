package update

import (
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/ccremer/greposync/infrastructure/ui"
	"github.com/go-logr/logr"
)

type updateInstrumentation struct {
	console *ui.ColoredConsole
	log     logr.Logger
}

func NewUpdateInstrumentation(console *ui.ColoredConsole, factory logging.LoggerFactory) *updateInstrumentation {
	return &updateInstrumentation{
		console: console,
		log:     factory.NewGenericLogger(""),
	}
}

func (i *updateInstrumentation) batchPipelineStarted(total int) {
	i.log.Info("Update started")
	p, _ := i.console.BatchProgressbar.WithTotal(total).Start()
	i.console.BatchProgressbar = p
	// TODO: determine quiet parameter via config
	i.console.Quiet = true
}

func (i *updateInstrumentation) batchPipelineCompleted() {
	i.log.Info("Update finished")
}

func (i *updateInstrumentation) pipelineForRepositoryStarted(repo *domain.GitRepository) {
	i.log.WithName(repo.URL.GetFullName()).V(logging.LevelDebug).Info("Starting pipeline")
}

func (i *updateInstrumentation) pipelineForRepositoryCompleted(repo *domain.GitRepository, err error) {
	i.console.PrintProgressbarMessage(repo.URL.GetFullName(), err)
}
