package update

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/ccremer/greposync/infrastructure/repositorystore"
	"github.com/ccremer/greposync/infrastructure/templateengine/gotemplate"
	"github.com/ccremer/greposync/infrastructure/ui"
)

type AppService struct {
	engine        *gotemplate.GoTemplateEngine
	repoStore     *repositorystore.RepositoryStore
	templateStore *gotemplate.GoTemplateStore
	valueStore    domain.ValueStore
	prStore       domain.PullRequestStore
	renderService *domain.RenderService
	diffPrinter   *ui.ConsoleDiffPrinter
	cfg           *cfg.Configuration
	sink          *ui.ConsoleSink
	console       *ui.ColoredConsole
}

func NewConfigurator(
	engine *gotemplate.GoTemplateEngine,
	repoStore *repositorystore.RepositoryStore,
	templateStore *gotemplate.GoTemplateStore,
	valueStore domain.ValueStore,
	prStore domain.PullRequestStore,
	renderService *domain.RenderService,
	diffPrinter *ui.ConsoleDiffPrinter,
	cfg *cfg.Configuration,
	sink *ui.ConsoleSink,
	console *ui.ColoredConsole,
) *AppService {
	return &AppService{
		engine:        engine,
		repoStore:     repoStore,
		templateStore: templateStore,
		valueStore:    valueStore,
		prStore:       prStore,
		renderService: renderService,
		diffPrinter:   diffPrinter,
		cfg:           cfg,
		sink:          sink,
		console:       console,
	}
}

func (c *AppService) ConfigureInfrastructure() {
	c.engine.RootDir = "template"
	c.repoStore.ParentDir = "repos"
	c.repoStore.DefaultNamespace = c.cfg.Git.Namespace
	c.repoStore.CommitBranch = c.cfg.Git.CommitBranch
	c.templateStore.RootDir = "template"
	c.console.Quiet = !c.cfg.Log.ShowLog
	level := logging.ParseLevelOrDefault(c.cfg.Log.Level, logging.LevelInfo)
	c.disableLogLevelsBelow(level)
}

func (c *AppService) disableLogLevelsBelow(level logging.LogLevel) {
	s := c.sink
	switch level {
	case logging.LevelWarn:
		s.
			WithLevelEnabled(logging.LevelSuccess, false).
			WithLevelEnabled(logging.LevelInfo, false).
			WithLevelEnabled(logging.LevelDebug, false)
	case logging.LevelSuccess:
		s.
			WithLevelEnabled(logging.LevelInfo, false).
			WithLevelEnabled(logging.LevelDebug, false)
	case logging.LevelInfo:
		s.
			WithLevelEnabled(logging.LevelDebug, false)
	}
}
