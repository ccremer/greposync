//+build wireinject

package main

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/cli"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/core/gitrepo"
	"github.com/ccremer/greposync/core/pullrequest"
	"github.com/ccremer/greposync/pkg/githosting/github"
	"github.com/ccremer/greposync/pkg/rendering"
	"github.com/ccremer/greposync/pkg/repository"
	"github.com/ccremer/greposync/pkg/valuestore"
	"github.com/google/wire"
)

type injector struct {
	prepareWorkspaceHandler *gitrepo.PrepareWorkspaceHandler
	pullRequestHandler      *pullrequest.PullRequestHandler
	app                     *cli.App
}

func NewInjector(
	app *cli.App,
	pwh *gitrepo.PrepareWorkspaceHandler,
	prh *pullrequest.PullRequestHandler,
) *injector {
	i := &injector{
		prepareWorkspaceHandler: pwh,
		pullRequestHandler:      prh,
		app:                     app,
	}
	return i
}

func (i *injector) RegisterHandlers() {
	core.RegisterHandler(gitrepo.PrepareWorkspaceEvent, i.prepareWorkspaceHandler)
	core.RegisterHandler(pullrequest.EnsurePullRequestEvent, i.pullRequestHandler)
}

func (i *injector) RunApp() {
	i.app.Run()
}

func initInjector() *injector {
	panic(wire.Build(
		NewInjector,

		// CLI
		cli.NewApp,
		wire.Value(cli.VersionInfo{Version: version, Commit: commit, Date: date}),
		cfg.NewDefaultConfig,

		// Core
		gitrepo.NewPrepareWorkspaceHandler,
		pullrequest.NewPullRequestHandler,

		// Stores
		wire.NewSet(repository.NewRepositoryStore, wire.Bind(new(core.GitRepositoryStore), new(*repository.RepositoryStore))),
		wire.NewSet(rendering.NewGoTemplateStore, wire.Bind(new(core.TemplateStore), new(*rendering.GoTemplateStore))),
		wire.NewSet(valuestore.NewValueStore, wire.Bind(new(core.ValueStore), new(*valuestore.KoanfValueStore))),

		// Git providers
		newGitProviders,
		github.NewRemote,
	))
}

func newGitProviders(ghRemote *github.GhRemote) repository.ProviderMap {
	return repository.ProviderMap{
		github.GitHubProviderKey: ghRemote,
	}
}
