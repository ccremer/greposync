package instrumentation

import (
	"github.com/ccremer/go-command-pipeline/parallel"
	"github.com/ccremer/greposync/domain"
)

type BatchInstrumentation interface {
	BatchPipelineStarted(repos []*domain.GitRepository)
	BatchPipelineCompleted(repos []*domain.GitRepository)
	PipelineForRepositoryStarted(repo *domain.GitRepository)
	PipelineForRepositoryCompleted(repo *domain.GitRepository, err error)
	NewCollectErrorHandler(skipBroken bool) parallel.ResultHandler
}

type InstrumentationContext interface {
	GetRepositories() []*domain.GitRepository
}
