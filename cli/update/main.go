package update

import (
	"fmt"
	"net/url"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/go-command-pipeline/parallel"
	"github.com/ccremer/go-command-pipeline/predicate"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/cli/flags"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/core/gitrepo"
	"github.com/ccremer/greposync/core/pullrequest"
	corerendering "github.com/ccremer/greposync/core/rendering"
	"github.com/ccremer/greposync/printer"
	"github.com/ccremer/greposync/repository"
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
		repositories []core.GitRepository
		repoStore    core.GitRepositoryStore
		globalK      *koanf.Koanf
	}
)

// NewCommand returns a new Command instance.
func NewCommand(cfg *cfg.Configuration, repoStore core.GitRepositoryStore) *Command {
	c := &Command{
		globalK:   koanf.New("."),
		cfg:       cfg,
		repoStore: repoStore,
	}
	c.cliCommand = c.createCliCommand()
	return c
}

// GetCliCommand returns the command instance for CLI library.
func (c *Command) GetCliCommand() *cli.Command {
	return c.cliCommand
}

func (c *Command) createCliCommand() *cli.Command {
	return &cli.Command{
		Name:   "update",
		Usage:  "Update the repositories in managed_repos.yml",
		Action: c.runCommand,
		Before: c.validateUpdateCommand,
		Flags: flags.CombineWithGlobalFlags(
			&cli.StringFlag{
				Name:    dryRunFlagName,
				Aliases: []string{"d"},
				Usage:   "Select a dry run mode. Allowed values: offline (do not run any Git commands), commit (commit, but don't push), push (push, but don't touch PRs)",
			},
			&cli.BoolFlag{
				Name:  amendFlagName,
				Usage: "Amend previous commit.",
			},
			&cli.BoolFlag{
				Name:  prCreateFlagName,
				Usage: "Create a PullRequest on a supported git hoster after pushing to remote.",
			},
			&cli.StringFlag{
				Name:  prBodyFlagName,
				Usage: "Markdown-enabled body of the PullRequest. It will load from an existing file if this is a path. Content can be templated. Defaults to commit message.",
			},
		),
	}
}

func (c *Command) runCommand(_ *cli.Context) error {

	logger := printer.PipelineLogger{Logger: printer.New().SetName("update").SetLevel(printer.DefaultLevel)}
	p := pipeline.NewPipelineWithLogger(logger).WithSteps(
		pipeline.NewStep("parse managed repos config", c.parseServices()),
		parallel.NewWorkerPoolStep("update repositories", c.cfg.Project.Jobs, c.updateReposInParallel(), c.errorHandler()),
	)
	return p.Run().Err
}

func (c *Command) createPipeline(r core.GitRepository) *pipeline.Pipeline {
	log := printer.New().SetName(r.GetConfig().URL.GetRepositoryName()).SetLevel(printer.DefaultLevel)

	sc := &cfg.SyncConfig{
		Git:         r.Config,
		PullRequest: c.cfg.PullRequest,
		Template: &cfg.TemplateConfig{
			RootDir: c.cfg.Template.RootDir,
		},
	}
	repoUrl := sc.Git.Url
	logger := printer.PipelineLogger{Logger: log}
	p := pipeline.NewPipelineWithLogger(logger)
	p.WithSteps(
		pipeline.NewStep("prepare workspace", c.prepareWorkspace(repoUrl)),
		pipeline.NewStep("render templates", c.renderTemplates(repoUrl)),
		predicate.WrapIn(pipeline.NewPipelineWithLogger(logger).
			WithNestedSteps("commit changes",
				pipeline.NewStep("add", r.Add()),
				pipeline.NewStep("commit", r.Commit()),
				pipeline.NewStep("show diff", r.Diff()),
			),
			predicate.And(r.EnabledCommit(), r.Dirty())),
		predicate.ToStep("push changes", r.PushToRemote(), predicate.And(r.EnabledPush(), r.IfBranchHasCommits())),
		predicate.ToStep("pull request", c.ensurePullRequest(repoUrl), predicate.And(r.IfBranchHasCommits(), predicate.Bool(sc.PullRequest.Create))),
		pipeline.NewStep("end", func() pipeline.Result {
			log.InfoF("Pipeline for '%s/%s' finished", sc.Git.Namespace, sc.Git.Name)
			return pipeline.Result{}
		}),
	)
	return p
}

func (c *Command) parseServices() func() pipeline.Result {
	return func() pipeline.Result {
		repos, err := c.repoStore.FetchGitRepositories()
		c.repositories = repos
		return pipeline.Result{Err: err}
	}
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
				err = multierror.Append(err, fmt.Errorf("%s: %w", repo.GetConfig().URL.GetRepositoryName(), result.Err))
			}
		}
		return pipeline.Result{Err: err}
	}
}

func (c *Command) prepareWorkspace(url *url.URL) pipeline.ActionFunc {
	return c.fireEvent(url, gitrepo.PrepareWorkspaceEvent)
}

func (c *Command) ensurePullRequest(url *url.URL) pipeline.ActionFunc {
	return c.fireEvent(url, pullrequest.EnsurePullRequestEvent)
}

func (c *Command) renderTemplates(url *url.URL) pipeline.ActionFunc {
	return c.fireEvent(url, corerendering.RenderTemplatesEvent)
}

func (c *Command) fireEvent(u *url.URL, event core.EventName) pipeline.ActionFunc {
	return func() pipeline.Result {
		result := <-core.FireEvent(event, core.EventSource{
			Url: core.FromURL(u),
		})
		return pipeline.Result{Err: result.Error}
	}
}
