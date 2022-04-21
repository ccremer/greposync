//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ccremer/greposync/application"
	"github.com/ccremer/greposync/application/initialize"
	"github.com/ccremer/greposync/application/instrumentation"
	"github.com/ccremer/greposync/application/labels"
	"github.com/ccremer/greposync/application/update"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/githosting"
	"github.com/ccremer/greposync/infrastructure/githosting/github"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/ccremer/greposync/infrastructure/repositorystore"
	"github.com/ccremer/greposync/infrastructure/templateengine"
	"github.com/ccremer/greposync/infrastructure/templateengine/gotemplate"
	"github.com/ccremer/greposync/infrastructure/ui"
	"github.com/ccremer/greposync/infrastructure/valuestore"
	"github.com/go-logr/logr"
	"github.com/google/wire"
)

type injector struct {
	app *application.App
}

func NewInjector(
	app *application.App,
) *injector {
	i := &injector{
		app: app,
	}
	return i
}

func (i *injector) RunApp() {
	i.app.Run()
}

func initInjector() *injector {
	panic(wire.Build(
		NewInjector,

		// CLI
		application.NewApp,
		update.NewConfigurator,
		wire.Value(application.VersionInfo{Version: version, Commit: commit, Date: date}),
		cfg.NewDefaultConfig,
		labels.NewCommand,
		labels.NewConfigurator,
		update.NewCommand,
		initialize.NewCommand,
		wire.NewSet(ui.NewConsoleDiffPrinter, wire.Bind(new(ui.DiffPrinter), new(*ui.ConsoleDiffPrinter))),

		// Template Engine
		wire.NewSet(gotemplate.NewEngine, wire.Bind(new(domain.TemplateEngine), new(*gotemplate.GoTemplateEngine))),
		wire.NewSet(gotemplate.NewTemplateStore, wire.Bind(new(domain.TemplateStore), new(*gotemplate.GoTemplateStore))),
		wire.NewSet(templateengine.NewRenderServiceInstrumentation, wire.Bind(new(domain.RenderServiceInstrumentation), new(*templateengine.RenderServiceInstrumentation))),
		wire.NewSet(templateengine.NewCleanupServiceInstrumentation, wire.Bind(new(domain.CleanupServiceInstrumentation), new(*templateengine.CleanupServiceInstrumentation))),

		// Stores
		wire.NewSet(repositorystore.NewRepositoryStore, wire.Bind(new(domain.GitRepositoryStore), new(*repositorystore.RepositoryStore))),
		wire.NewSet(valuestore.NewKoanfStore, wire.Bind(new(domain.ValueStore), new(*valuestore.KoanfStore))),
		wire.NewSet(githosting.NewPullRequestStore, wire.Bind(new(domain.PullRequestStore), new(*githosting.PullRequestStore))),
		wire.NewSet(githosting.NewLabelStore, wire.Bind(new(domain.LabelStore), new(*githosting.LabelStore))),

		// Services
		domain.NewRenderService,
		domain.NewCleanupService,
		domain.NewPullRequestService,

		// Console & Logging
		wire.NewSet(ui.NewConsoleSink, wire.Bind(new(logr.LogSink), new(*ui.ConsoleSink))),
		wire.NewSet(ui.NewConsoleLoggerFactory, wire.Bind(new(logging.LoggerFactory), new(*ui.ConsoleLoggerFactory))),
		ui.NewColoredConsole,

		// Instrumentation
		valuestore.NewValueStoreInstrumentation,
		wire.NewSet(instrumentation.NewUpdateInstrumentation, wire.Bind(new(instrumentation.BatchInstrumentation), new(*instrumentation.CommonBatchInstrumentation))),
		repositorystore.NewRepositoryStoreInstrumentation,
		github.NewGitHubInstrumentation,

		// Git providers
		newGitProviders,
		github.NewRemote,
	))
}

func newGitProviders(ghRemote *github.GhRemote) githosting.ProviderMap {
	return githosting.ProviderMap{
		github.ProviderKey: ghRemote,
	}
}
