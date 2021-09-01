package update

import (
	"fmt"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/go-command-pipeline/parallel"
	"github.com/ccremer/go-command-pipeline/predicate"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/printer"
	"github.com/hashicorp/go-multierror"
	"github.com/knadh/koanf"
	"github.com/urfave/cli/v2"
)

const (
	dryRunFlagName   = "dry-run"
	prCreateFlagName = "pr-create"
	prBodyFlagName   = "pr-bodyTemplate"
	amendFlagName    = "git-amend"
)

type (
	// Command is a facade service for the update command that holds all dependent services and settings.
	Command struct {
		cfg          *cfg.Configuration
		cliCommand   *cli.Command
		repositories []*domain.GitRepository
		globalK      *koanf.Koanf
		appService   *AppService
	}
)

// NewCommand returns a new Command instance.
func NewCommand(
	cfg *cfg.Configuration,
	configurator *AppService,
) *Command {
	c := &Command{
		globalK:    koanf.New("."),
		cfg:        cfg,
		appService: configurator,
	}
	return c
}

func (c *Command) runCommand(_ *cli.Context) error {

	logger := printer.PipelineLogger{Logger: printer.New().SetName("update").SetLevel(printer.DefaultLevel)}
	p := pipeline.NewPipeline().AddBeforeHook(logger).WithSteps(
		pipeline.NewStep("configure infrastructure", c.configureInfrastructure()),
		pipeline.NewStep("fetch managed repos config", c.fetchRepositories()),
		parallel.NewWorkerPoolStep("update repositories", c.cfg.Project.Jobs, c.updateReposInParallel(), c.errorHandler()),
	)
	return p.Run().Err
}

func (c *Command) createPipeline(r *domain.GitRepository) *pipeline.Pipeline {
	log := printer.New().SetName(r.URL.GetFullName()).SetLevel(printer.DefaultLevel)

	sc := &cfg.SyncConfig{
		PullRequest: c.cfg.PullRequest,
		Template: &cfg.TemplateConfig{
			RootDir: c.cfg.Template.RootDir,
		},
	}

	resetRepo := !c.cfg.Git.SkipReset
	enabledCommits := !c.cfg.Git.SkipCommit
	enabledPush := !c.cfg.Git.SkipPush

	repoCtx := &pipelineContext{
		repo:       r,
		appService: c.appService,
		differ: &Differ{
			log:        printer.New().MapColorToLevel(printer.Blue, printer.LevelInfo).SetName(r.URL.GetRepositoryName()),
			repository: r,
		},
	}

	logger := printer.PipelineLogger{Logger: log}
	p := pipeline.NewPipeline().AddBeforeHook(logger)
	p.WithSteps(
		pipeline.NewPipeline().AddBeforeHook(logger).
			WithNestedSteps("prepare workspace",
				predicate.ToStep("clone repository", repoCtx.clone(), repoCtx.dirMissing()),
				predicate.ToStep("fetch", repoCtx.fetch(), predicate.Bool(resetRepo)),
				predicate.ToStep("reset", repoCtx.reset(), predicate.Bool(resetRepo)),
				predicate.ToStep("checkout branch", repoCtx.checkout(), predicate.Bool(resetRepo)),
				predicate.ToStep("pull", repoCtx.fetch(), predicate.Bool(resetRepo)),
			),
		pipeline.NewStep("render templates", repoCtx.renderTemplates()),
		predicate.WrapIn(pipeline.NewPipeline().AddBeforeHook(logger).
			WithNestedSteps("commit changes",
				pipeline.NewStep("add", repoCtx.add()),
				pipeline.NewStep("commit", repoCtx.commit()),
				pipeline.NewStep("show diff", repoCtx.diff()),
			),
			predicate.And(predicate.Bool(enabledCommits), repoCtx.isDirty())),
		predicate.ToStep("push changes", repoCtx.push(), predicate.And(predicate.Bool(enabledPush), repoCtx.hasCommits())),
		predicate.ToStep("find existing pull request", repoCtx.fetchPullRequest(), predicate.Bool(sc.PullRequest.Create)),
		predicate.ToStep("update pull request", repoCtx.ensurePullRequest(), predicate.And(repoCtx.hasCommits(), predicate.Bool(sc.PullRequest.Create))),
		pipeline.NewStep("end", func(_ pipeline.Context) pipeline.Result {
			log.InfoF("Pipeline for '%s/%s' finished", r.URL.GetNamespace(), r.URL.GetRepositoryName())
			return pipeline.Result{}
		}),
	)
	return p
}

func (c *Command) updateReposInParallel() parallel.PipelineSupplier {
	return func(pipelines chan *pipeline.Pipeline) {
		defer close(pipelines)
		for _, r := range c.repositories {
			p := c.createPipeline(r)
			pipelines <- p
		}
	}
}

func (c *Command) errorHandler() parallel.ResultHandler {
	return func(results map[uint64]pipeline.Result) pipeline.Result {
		var err error
		for index, repo := range c.repositories {
			if result := results[uint64(index)]; result.Err != nil {
				err = multierror.Append(err, fmt.Errorf("%s: %w", repo.URL.GetRepositoryName(), result.Err))
			}
		}
		return pipeline.Result{Err: err}
	}
}

func (c *Command) configureInfrastructure() pipeline.ActionFunc {
	return func(_ pipeline.Context) pipeline.Result {
		c.appService.ConfigureInfrastructure()
		return pipeline.Result{}
	}
}

func (c *Command) fetchRepositories() pipeline.ActionFunc {
	return func(_ pipeline.Context) pipeline.Result {
		repos, err := c.appService.repoStore.FetchGitRepositories()
		c.repositories = repos
		return pipeline.Result{Err: err}
	}
}
