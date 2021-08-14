//+build wireinject

package main

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/cli"
	"github.com/ccremer/greposync/cli/labels"
	"github.com/ccremer/greposync/cli/update"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/core/gitrepo"
	corelabels "github.com/ccremer/greposync/core/labels"
	"github.com/ccremer/greposync/core/pullrequest"
	corerendering "github.com/ccremer/greposync/core/rendering"
	"github.com/ccremer/greposync/pkg/githosting/github"
	"github.com/ccremer/greposync/pkg/rendering"
	"github.com/ccremer/greposync/pkg/repository"
	"github.com/ccremer/greposync/pkg/valuestore"
	"github.com/google/wire"
)

type injector struct {
	app      *cli.App
	handlers map[core.EventName]core.EventHandler
}

func NewInjector(
	app *cli.App,
	pwh *gitrepo.PrepareWorkspaceHandler,
	prh *pullrequest.PullRequestHandler,
	rth *corerendering.RenderTemplatesHandler,
	luh *corelabels.LabelUpdateHandler,
) *injector {
	i := &injector{
		app:      app,
		handlers: map[core.EventName]core.EventHandler{},
	}
	i.handlers[gitrepo.PrepareWorkspaceEvent] = pwh
	i.handlers[pullrequest.EnsurePullRequestEvent] = prh
	i.handlers[corerendering.RenderTemplatesEvent] = rth
	i.handlers[corelabels.LabelUpdateEvent] = luh
	return i
}

func (i *injector) RegisterHandlers() {
	for name, handler := range i.handlers {
		core.RegisterHandler(name, handler)
	}
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
		labels.NewCommand,
		update.NewCommand,

		// Core
		gitrepo.NewPrepareWorkspaceHandler,
		pullrequest.NewPullRequestHandler,
		corerendering.NewRenderTemplatesHandler,
		corelabels.NewLabelUpdateHandler,

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
