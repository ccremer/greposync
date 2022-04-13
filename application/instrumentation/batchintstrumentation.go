package instrumentation

import (
	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/domain"
)

type BatchInstrumentation interface {
	BatchPipelineStarted(repos []*domain.GitRepository)
	BatchPipelineCompleted(repos []*domain.GitRepository)
	PipelineForRepositoryStarted(repo *domain.GitRepository)
	PipelineForRepositoryCompleted(repo *domain.GitRepository, err error)
	NewCollectErrorHandler(skipBroken bool) pipeline.ParallelResultHandler
}

type RepositoriesContextKey struct{}
