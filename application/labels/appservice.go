package labels

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/ccremer/greposync/infrastructure/repositorystore"
)

type AppService struct {
	repoStore  *repositorystore.RepositoryStore
	labelStore domain.LabelStore
	cfg        *cfg.Configuration
	factory    logging.LoggerFactory
}

func NewConfigurator(
	repoStore *repositorystore.RepositoryStore,
	labelStore domain.LabelStore,
	cfg *cfg.Configuration,
	factory logging.LoggerFactory,
) *AppService {
	return &AppService{
		repoStore:  repoStore,
		labelStore: labelStore,
		cfg:        cfg,
		factory:    factory,
	}
}
