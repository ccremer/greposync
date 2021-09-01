// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/ccremer/greposync/application"
	"github.com/ccremer/greposync/application/labels"
	"github.com/ccremer/greposync/application/update"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/githosting"
	"github.com/ccremer/greposync/infrastructure/githosting/github"
	"github.com/ccremer/greposync/infrastructure/repositorystore"
	"github.com/ccremer/greposync/infrastructure/templateengine"
	"github.com/ccremer/greposync/infrastructure/templateengine/gotemplate"
	"github.com/ccremer/greposync/infrastructure/valuestore"
)

// Injectors from wire.go:

func initInjector() *injector {
	versionInfo := _wireVersionInfoValue
	configuration := cfg.NewDefaultConfig()
	repositoryStore := repositorystore.NewRepositoryStore()
	ghRemote := github.NewRemote()
	providerMap := newGitProviders(ghRemote)
	labelStore := githosting.NewLabelStore(providerMap)
	appService := labels.NewConfigurator(repositoryStore, labelStore, configuration)
	command := labels.NewCommand(configuration, appService)
	goTemplateEngine := gotemplate.NewEngine()
	goTemplateStore := gotemplate.NewTemplateStore()
	koanfValueStore := valuestore.NewValueStore()
	pullRequestStore := githosting.NewPullRequestStore(providerMap)
	renderServiceInstrumentation := templateengine.NewRenderServiceInstrumentation()
	renderService := domain.NewRenderService(renderServiceInstrumentation)
	updateAppService := update.NewConfigurator(goTemplateEngine, repositoryStore, goTemplateStore, koanfValueStore, pullRequestStore, renderService, configuration)
	updateCommand := update.NewCommand(configuration, updateAppService)
	app := application.NewApp(versionInfo, configuration, command, updateCommand)
	mainInjector := NewInjector(app)
	return mainInjector
}

var (
	_wireVersionInfoValue = application.VersionInfo{Version: version, Commit: commit, Date: date}
)

// wire.go:

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

func newGitProviders(ghRemote *github.GhRemote) githosting.ProviderMap {
	return githosting.ProviderMap{github.ProviderKey: ghRemote}
}
