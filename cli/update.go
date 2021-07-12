package cli

import (
	"encoding/json"
	"fmt"
	"os"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/go-command-pipeline/parallel"
	"github.com/ccremer/go-command-pipeline/predicate"
	"github.com/ccremer/greposync/cfg"
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
	prBodyFlagName   = "pr-body"
	amendFlagName    = "git-amend"
)

var (
	JobsMinimumCount = 1
	JobsMaximumCount = 8
)

type (
	// UpdateCommand is a facade service for the update command that holds all dependent services and settings.
	UpdateCommand struct {
		cfg          *cfg.Configuration
		cliCommand   *cli.Command
		repoServices []*repository.Service
		parser       *rendering.Parser
		globalK      *koanf.Koanf
	}
)

// NewUpdateCommand returns a new UpdateCommand instance.
func NewUpdateCommand(cfg *cfg.Configuration) *UpdateCommand {
	return &UpdateCommand{
		globalK: koanf.New("."),
		cfg:     cfg,
		parser:  rendering.NewParser(cfg.Template),
	}
}

func (c *UpdateCommand) createUpdateCommand() *cli.Command {
	c.cliCommand = &cli.Command{
		Name:   "update",
		Usage:  "Update the repositories in managed_repos.yml",
		Action: c.runUpdateCommand,
		Before: c.validateUpdateCommand,
		Flags: combineWithGlobalFlags(
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
	return c.cliCommand
}

func (c *UpdateCommand) validateUpdateCommand(ctx *cli.Context) error {
	if err := cfg.ParseConfig(GrepoSyncFileName, config, ctx); err != nil {
		return err
	}

	if err := validateGlobalFlags(ctx); err != nil {
		return err
	}

	if jobs := config.Project.Jobs; jobs > JobsMaximumCount || jobs < JobsMinimumCount {
		return fmt.Errorf("--%s is required to be between %d and %d", projectJobsFlagName, JobsMinimumCount, JobsMaximumCount)
	}

	if ctx.IsSet(dryRunFlagName) {
		dryRunMode := ctx.String(dryRunFlagName)
		switch dryRunMode {
		case "offline":
			config.Git.SkipReset = true
			config.Git.SkipCommit = true
			config.Git.SkipPush = true
			config.PullRequest.Create = false
		case "commit":
			config.Git.SkipPush = true
			config.PullRequest.Create = false
		case "push":
			config.PullRequest.Create = false
		default:
			return fmt.Errorf("invalid flag value of %s: %s", dryRunFlagName, dryRunMode)
		}
	}

	config.Sanitize()
	j, _ := json.Marshal(config)
	printer.DebugF("Using config: %s", j)
	return nil
}

func (c *UpdateCommand) runUpdateCommand(_ *cli.Context) error {

	logger := printer.PipelineLogger{Logger: printer.New().SetName("update").SetLevel(printer.DefaultLevel)}
	p := pipeline.NewPipelineWithLogger(logger).WithSteps(
		pipeline.NewStep("parse config defaults", c.loadGlobalDefaults()),
		pipeline.NewStep("parse templates", c.parser.ParseTemplateDirAction()),
		pipeline.NewStep("parse managed repos config", c.parseServices()),
		parallel.NewWorkerPoolStep("update repositories", config.Project.Jobs, c.updateReposInParallel(), c.errorHandler()),
	)
	return p.Run().Err
}

func (c *UpdateCommand) createPipeline(r *repository.Service) *pipeline.Pipeline {
	log := printer.New().SetName(r.Config.Name).SetLevel(printer.DefaultLevel)

	sc := &cfg.SyncConfig{
		Git:         r.Config,
		PullRequest: config.PullRequest,
		Template: &cfg.TemplateConfig{
			RootDir: config.Template.RootDir,
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
		predicate.WrapIn(pipeline.NewPipelineWithLogger(logger).WithSteps(
			pipeline.NewStep("add", r.Add()),
			pipeline.NewStep("commit", r.Commit()),
			pipeline.NewStep("show diff", r.Diff()),
			predicate.ToStep("push", r.PushToRemote(), r.EnabledPush()),
		).AsNestedStep("push changes"), predicate.And(r.EnabledCommit(), r.Dirty())),
		predicate.WrapIn(pipeline.NewPipelineWithLogger(logger).WithSteps(
			pipeline.NewStep("render pull request template", renderer.RenderPrTemplate()),
			pipeline.NewStep("create or update pull request", r.CreateOrUpdatePr(config.PullRequest)),
		).AsNestedStep("pull request"), predicate.Bool(sc.PullRequest.Create)),
		pipeline.NewStep("end", func() pipeline.Result {
			log.InfoF("Pipeline for '%s/%s' finished", sc.Git.Namespace, sc.Git.Name)
			return pipeline.Result{}
		}),
	)
	return p
}

func (c *UpdateCommand) parseServices() func() pipeline.Result {
	return func() pipeline.Result {
		c.repoServices = repository.NewServicesFromFile(config)
		return pipeline.Result{}
	}
}

func (c *UpdateCommand) loadGlobalDefaults() pipeline.ActionFunc {
	return func() pipeline.Result {
		if info, err := os.Stat(ConfigDefaultName); err != nil || info.IsDir() {
			printer.WarnF("File %s does not exist, ignoring template defaults")
			return pipeline.Result{}
		}
		printer.DebugF("Loading config %s", ConfigDefaultName)
		err := c.globalK.Load(file.Provider(ConfigDefaultName), yaml.Parser())
		return pipeline.Result{Err: err}
	}
}

func (c *UpdateCommand) updateReposInParallel() parallel.PipelineSupplier {
	return func(pipelines chan *pipeline.Pipeline) {
		defer close(pipelines)
		for _, r := range c.repoServices {
			p := c.createPipeline(r)
			pipelines <- p
		}
	}
}

func (c *UpdateCommand) errorHandler() parallel.ResultHandler {
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
