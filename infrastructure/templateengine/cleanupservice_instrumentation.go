package templateengine

import (
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/go-logr/logr"
)

type CleanupServiceInstrumentation struct {
	factory logging.LoggerFactory
	log     logr.Logger
}

func NewCleanupServiceInstrumentation(factory logging.LoggerFactory) *CleanupServiceInstrumentation {
	i := &CleanupServiceInstrumentation{
		factory: factory,
	}
	return i
}

func (r *CleanupServiceInstrumentation) WithRepository(repository *domain.GitRepository) domain.CleanupServiceInstrumentation {
	newCopy := NewCleanupServiceInstrumentation(r.factory)
	newCopy.log = r.factory.NewRepositoryLogger(repository)
	return newCopy
}

func (r *CleanupServiceInstrumentation) FetchedFilesToDelete(fetchErr error, files []domain.Path) error {
	if fetchErr == nil {
		r.log.V(logging.LevelDebug).Info("Fetched files", "files", files)
	}
	return fetchErr
}

func (r *CleanupServiceInstrumentation) DeletedFile(file domain.Path) {
	r.log.V(logging.LevelDebug).Info("Deleted file", "file", file)
}
