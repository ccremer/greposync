package update

import (
	"fmt"
	"os"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/go-command-pipeline/parallel"
	"github.com/ccremer/go-command-pipeline/predicate"
	"github.com/ccremer/greposync/cfg"
	app "github.com/ccremer/greposync/cli"
	"github.com/ccremer/greposync/printer"
	"github.com/ccremer/greposync/rendering"
	"github.com/ccremer/greposync/repository"
	"github.com/hashicorp/go-multierror"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
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
		repoServices []*repository.Service
		parser       *rendering.Parser
		globalK      *koanf.Koanf
	}
)

// NewCommand returns a new Command instance.
func NewCommand(cfg *cfg.Configuration) *Command {
	c := &Command{
		globalK: koanf.New("."),
		cfg:     cfg,
		parser:  rendering.NewParser(cfg.Template),
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
		Flags: app.CombineWithGlobalFlags(
			app.NewProjectIncludeFlag(),
			app.NewProjectExcludeFlag(),
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
		pipeline.NewStep("parse config defaults", c.loadGlobalDefaults()),
		pipeline.NewStep("parse templates", c.parser.ParseTemplateDirAction()),
		pipeline.NewStep("parse managed repos config", c.parseServices()),
		parallel.NewWorkerPoolStep("update repositories", c.cfg.Project.Jobs, c.updateReposInParallel(), c.errorHandler()),
	)
	return p.Run().Err
}

func (c *Command) createPipeline(r *repository.Service) *pipeline.Pipeline {
	log := printer.New().SetName(r.Config.Name).SetLevel(printer.DefaultLevel)

	sc := &cfg.SyncConfig{
		Git:         r.Config,
		PullRequest: c.cfg.PullRequest,
		Template: &cfg.TemplateConfig{
			RootDir: c.cfg.Template.RootDir,
		},
	}
	renderer := rendering.NewRenderer(sc, c.globalK, c.parser)
	gitDirExists := r.DirExists(r.Config.Dir)
	logger := printer.PipelineLogger{Logger: log}
	p := pipeline.NewPipelineWithLogger(logger)
	p.WithSteps(
		pipeline.NewPipelineWithLogger(logger).WithSteps(
			predicate.ToStep("clone repository", r.CloneGitRepository(), predicate.Bool(!gitDirExists)),
			pipeline.NewStep("determine default branch", r.GetDefaultBranch()),
			predicate.ToStep("fetch", r.Fetch(), r.EnabledReset()),
			predicate.ToStep("reset repository", r.ResetRepository(), r.EnabledReset()),
			pipeline.NewStep("checkout branch", r.CheckoutBranch()),
			predicate.ToStep("pull", r.Pull(), r.EnabledReset()),
		).AsNestedStep("prepare workspace"),
		pipeline.NewStep("render templates", renderer.RenderTemplateDir()),
		pipeline.NewStep("cleanup unwanted files", renderer.DeleteUnwantedFiles()),
		predicate.WrapIn(pipeline.NewPipelineWithLogger(logger).WithSteps(
			pipeline.NewStep("add", r.Add()),
			pipeline.NewStep("commit", r.Commit()),
			pipeline.NewStep("show diff", r.Diff()),
		).AsNestedStep("commit changes"), predicate.And(r.EnabledCommit(), r.Dirty())),
		predicate.ToStep("push changes", r.PushToRemote(), predicate.And(r.EnabledPush(), r.IfBranchHasCommits())),
		predicate.WrapIn(pipeline.NewPipelineWithLogger(logger).WithSteps(
			pipeline.NewStep("render pull request template", renderer.RenderPrTemplate()),
			pipeline.NewStep("prepare API", r.InitializeGitHubProvider(c.cfg.PullRequest)),
			pipeline.NewStep("create or update pull request", r.CreateOrUpdatePr(c.cfg.PullRequest)),
		).AsNestedStep("pull request"), predicate.And(r.IfBranchHasCommits(), predicate.Bool(sc.PullRequest.Create))),
		pipeline.NewStep("end", func() pipeline.Result {
			log.InfoF("Pipeline for '%s/%s' finished", sc.Git.Namespace, sc.Git.Name)
			return pipeline.Result{}
		}),
	)
	return p
}

func (c *Command) parseServices() func() pipeline.Result {
	return func() pipeline.Result {
		s, err := repository.NewServicesFromFile(c.cfg)
		c.repoServices = s
		return pipeline.Result{Err: err}
	}
}

func (c *Command) loadGlobalDefaults() pipeline.ActionFunc {
	return func() pipeline.Result {
		if info, err := os.Stat(app.ConfigDefaultName); err != nil || info.IsDir() {
			printer.WarnF("File %s does not exist, ignoring template defaults")
			return pipeline.Result{}
		}
		printer.DebugF("Loading config %s", app.ConfigDefaultName)
		err := c.globalK.Load(file.Provider(app.ConfigDefaultName), yaml.Parser())
		return pipeline.Result{Err: err}
	}
}

func (c *Command) updateReposInParallel() parallel.PipelineSupplier {
	return func(pipelines chan *pipeline.Pipeline) {
		defer close(pipelines)
		for _, r := range c.repoServices {
			p := c.createPipeline(r)
			pipelines <- p
		}
	}
}

func (c *Command) errorHandler() parallel.ResultHandler {
	return func(results map[uint64]pipeline.Result) pipeline.Result {
		var err error
		for index, service := range c.repoServices {
			if result := results[uint64(index)]; result.Err != nil {
				err = multierror.Append(err, fmt.Errorf("%s: %w", service.Config.Name, result.Err))
			}
		}
		return pipeline.Result{Err: err}
	}
}
