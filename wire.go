//+build wireinject

package main

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/cli"
	"github.com/ccremer/greposync/cli/labels"
	"github.com/ccremer/greposync/cli/update"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/githosting"
	"github.com/ccremer/greposync/infrastructure/githosting/github"
	"github.com/ccremer/greposync/infrastructure/repositorystore"
	"github.com/ccremer/greposync/infrastructure/templateengine/gotemplate"
	"github.com/ccremer/greposync/infrastructure/valuestore"
	"github.com/ccremer/greposync/pkg/repository"
	"github.com/google/wire"
)

type injector struct {
	app *cli.App
}

func NewInjector(
	app *cli.App,
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
		cli.NewApp,
		update.NewConfigurator,
		wire.Value(cli.VersionInfo{Version: version, Commit: commit, Date: date}),
		cfg.NewDefaultConfig,
		labels.NewCommand,
		update.NewCommand,

		// Template Engine
		wire.NewSet(gotemplate.NewEngine, wire.Bind(new(domain.TemplateEngine), new(*gotemplate.GoTemplateEngine))),
		wire.NewSet(gotemplate.NewTemplateStore, wire.Bind(new(domain.TemplateStore), new(*gotemplate.GoTemplateStore))),

		// Stores
		wire.NewSet(repository.NewRepositoryStore, wire.Bind(new(core.GitRepositoryStore), new(*repository.RepositoryStore))),
		//wire.NewSet(rendering.NewGoTemplateStore, wire.Bind(new(core.GoTemplateStore), new(*rendering.GoTemplateStore))),
		wire.NewSet(valuestore.NewValueStore, wire.Bind(new(domain.ValueStore), new(*valuestore.KoanfValueStore))),
		wire.NewSet(githosting.NewPullRequestStore, wire.Bind(new(domain.PullRequestStore), new(*githosting.PullRequestStore))),
		wire.NewSet(githosting.NewLabelStore, wire.Bind(new(domain.LabelStore), new(*githosting.LabelStore))),
		wire.NewSet(repositorystore.NewRepositoryStore, wire.Bind(new(domain.GitRepositoryStore), new(*repositorystore.RepositoryStore))),

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
