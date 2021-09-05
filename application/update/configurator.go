package update

import (
	"github.com/ccremer/greposync/cfg"
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
	cfg           *cfg.Configuration
}

func NewConfigurator(
	engine *gotemplate.GoTemplateEngine,
	repoStore *repositorystore.RepositoryStore,
	templateStore *gotemplate.GoTemplateStore,
	valueStore domain.ValueStore,
	prStore domain.PullRequestStore,
	renderService *domain.RenderService,
	cfg *cfg.Configuration,
) *AppService {
	return &AppService{
		engine:        engine,
		repoStore:     repoStore,
		templateStore: templateStore,
		valueStore:    valueStore,
		prStore:       prStore,
		renderService: renderService,
		cfg:           cfg,
	}
}

func (c *AppService) ConfigureInfrastructure() {
	c.engine.RootDir = "template"
	c.repoStore.ParentDir = "repos"
	c.repoStore.DefaultNamespace = c.cfg.Git.Namespace
	c.repoStore.CommitBranch = c.cfg.Git.CommitBranch
	c.templateStore.RootDir = "template"
}
