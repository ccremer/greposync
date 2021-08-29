package labels

import (
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/repositorystore"
)

type AppService struct {
	repoStore  *repositorystore.RepositoryStore
	labelStore domain.LabelStore
}

func NewConfigurator(
	repoStore *repositorystore.RepositoryStore,
	labelStore domain.LabelStore,
) *AppService {
	return &AppService{
		repoStore:  repoStore,
		labelStore: labelStore,
	}
}

func (c *AppService) ConfigureInfrastructure() {
	c.repoStore.ParentDir = "repos"
}
