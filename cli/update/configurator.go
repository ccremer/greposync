package update

import (
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/repositorystore"
	"github.com/ccremer/greposync/infrastructure/templateengine/gotemplate"
)

type AppService struct {
	engine        *gotemplate.GoTemplateEngine
	repoStore     *repositorystore.RepositoryStore
	templateStore *gotemplate.GoTemplateStore
	valueStore    domain.ValueStore
}

func NewConfigurator(
	engine *gotemplate.GoTemplateEngine,
	repoStore *repositorystore.RepositoryStore,
	templateStore *gotemplate.GoTemplateStore,
	valueStore domain.ValueStore,
) *AppService {
	return &AppService{
		engine:        engine,
		repoStore:     repoStore,
		templateStore: templateStore,
		valueStore:    valueStore,
	}
}

func (c *AppService) ConfigureInfrastructure() {
	c.engine.RootDir = "template"
	c.repoStore.ParentDir = "repos"
	c.templateStore.RootDir = "template"
}
