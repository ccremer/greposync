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
	prStore       domain.PullRequestStore
	renderService *domain.RenderService
}

func NewConfigurator(
	engine *gotemplate.GoTemplateEngine,
	repoStore *repositorystore.RepositoryStore,
	templateStore *gotemplate.GoTemplateStore,
	valueStore domain.ValueStore,
	prStore domain.PullRequestStore,
	renderService *domain.RenderService,
) *AppService {
	return &AppService{
		engine:        engine,
		repoStore:     repoStore,
		templateStore: templateStore,
		valueStore:    valueStore,
		prStore:       prStore,
		renderService: renderService,
	}
}

func (c *AppService) ConfigureInfrastructure() {
	c.engine.RootDir = "template"
	c.repoStore.ParentDir = "repos"
	c.templateStore.RootDir = "template"
}
