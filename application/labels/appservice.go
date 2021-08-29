package labels

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/repositorystore"
)

type AppService struct {
	repoStore  *repositorystore.RepositoryStore
	labelStore domain.LabelStore
	cfg        *cfg.Configuration
}

func NewConfigurator(
	repoStore *repositorystore.RepositoryStore,
	labelStore domain.LabelStore,
	cfg *cfg.Configuration,
) *AppService {
	return &AppService{
		repoStore:  repoStore,
		labelStore: labelStore,
		cfg:        cfg,
	}
}

func (c *AppService) ConfigureInfrastructure() {
	c.repoStore.ParentDir = "repos"
	c.repoStore.DefaultNamespace = c.cfg.Git.Namespace
}
