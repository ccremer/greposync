package test

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/repositorystore"
	"github.com/ccremer/greposync/infrastructure/templateengine/gotemplate"
	"github.com/ccremer/greposync/infrastructure/ui"
)

type AppService struct {
	engine         *gotemplate.GoTemplateEngine
	repoStore      *repositorystore.TestRepositoryStore
	templateStore  *gotemplate.GoTemplateStore
	valueStore     domain.ValueStore
	renderService  *domain.RenderService
	cleanupService *domain.CleanupService
	diffPrinter    *ui.ConsoleDiffPrinter
	cfg            *cfg.Configuration
	console        *ui.ColoredConsole
}

func NewConfigurator(
	engine *gotemplate.GoTemplateEngine,
	repoStore *repositorystore.TestRepositoryStore,
	templateStore *gotemplate.GoTemplateStore,
	valueStore domain.ValueStore,
	renderService *domain.RenderService,
	cleanupService *domain.CleanupService,
	diffPrinter *ui.ConsoleDiffPrinter,
	cfg *cfg.Configuration,
	console *ui.ColoredConsole,
) *AppService {
	return &AppService{
		engine:         engine,
		repoStore:      repoStore,
		templateStore:  templateStore,
		valueStore:     valueStore,
		renderService:  renderService,
		cleanupService: cleanupService,
		diffPrinter:    diffPrinter,
		cfg:            cfg,
		console:        console,
	}
}

func (c *AppService) ConfigureInfrastructure() {
	c.engine.RootDir = "template"
	c.repoStore.ParentDir = "tests"
	c.repoStore.TestOutputRootDir = ".tests"
	c.repoStore.DefaultNamespace = "local"
	c.repoStore.ExcludeFilter = c.cfg.Project.Exclude
	c.repoStore.IncludeFilter = c.cfg.Project.Include
	c.templateStore.RootDir = "template"
	c.templateStore.SkipRemovingFileExtension = true
	c.console.SetTitle("RUNNING TESTS...")
	c.console.SetCommandName("Test")
}
