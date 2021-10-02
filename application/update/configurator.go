package update

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/repositorystore"
	"github.com/ccremer/greposync/infrastructure/templateengine/gotemplate"
	"github.com/ccremer/greposync/infrastructure/ui"
)

type AppService struct {
	engine         *gotemplate.GoTemplateEngine
	repoStore      *repositorystore.RepositoryStore
	templateStore  *gotemplate.GoTemplateStore
	valueStore     domain.ValueStore
	prStore        domain.PullRequestStore
	renderService  *domain.RenderService
	diffPrinter    *ui.ConsoleDiffPrinter
	cfg            *cfg.Configuration
	console        *ui.ColoredConsole
	cleanupService *domain.CleanupService
	prService      *domain.PullRequestService
}

func NewConfigurator(
	engine *gotemplate.GoTemplateEngine,
	repoStore *repositorystore.RepositoryStore,
	templateStore *gotemplate.GoTemplateStore,
	valueStore domain.ValueStore,
	prStore domain.PullRequestStore,
	renderService *domain.RenderService,
	cleanupService *domain.CleanupService,
	pullRequestService *domain.PullRequestService,
	diffPrinter *ui.ConsoleDiffPrinter,
	cfg *cfg.Configuration,
	console *ui.ColoredConsole,
) *AppService {
	return &AppService{
		engine:         engine,
		repoStore:      repoStore,
		templateStore:  templateStore,
		valueStore:     valueStore,
		prStore:        prStore,
		renderService:  renderService,
		cleanupService: cleanupService,
		prService:      pullRequestService,
		diffPrinter:    diffPrinter,
		cfg:            cfg,
		console:        console,
	}
}

func (c *AppService) ConfigureInfrastructure() {
	c.engine.RootDir = "template"
	c.repoStore.ParentDir = "repos"
	c.repoStore.DefaultNamespace = c.cfg.Git.Namespace
	c.repoStore.CommitBranch = c.cfg.Git.CommitBranch
	c.templateStore.RootDir = "template"
	c.console.Quiet = !c.cfg.Log.ShowLog
}
