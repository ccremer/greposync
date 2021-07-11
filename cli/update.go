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
	createPrFlagName = "pr"
	prBodyFlagName   = "pr-body"
	amendFlagName    = "amend"
)

func createUpdateCommand(c *cfg.Configuration) *cli.Command {
	return &cli.Command{
		Name:   "update",
		Usage:  "Update the repositories in managed_repos.yml",
		Action: runUpdateCommand,
		Before: validateUpdateCommand,
		Flags: combineWithGlobalFlags(
			&cli.StringFlag{
				Name:    dryRunFlagName,
				Aliases: []string{"d"},
				Usage:   "Select a dry run mode. Allowed values: offline (do not run any Git commands), commit (commit, but don't push), push (push, but don't touch PRs)",
			},
			&cli.BoolFlag{
				Name:        createPrFlagName,
				Destination: &c.PullRequest.Create,
				Usage:       "Create a PullRequest on a supported git hoster after pushing to remote.",
			},
			&cli.BoolFlag{
				Name:        amendFlagName,
				Destination: &c.Git.Amend,
				Usage:       "Amend previous commit.",
			},
			&cli.StringFlag{
				Name:  prBodyFlagName,
				Usage: "Markdown-enabled body of the PullRequest. It will load from an existing file if this is a path. Content can be templated. Defaults to commit message.",
			},
		),
	}
}

func validateUpdateCommand(ctx *cli.Context) error {
	if err := cfg.ParseConfig("greposync.yml", config); err != nil {
		return err
	}

	if err := validateGlobalFlags(ctx); err != nil {
		return err
	}

	if ctx.Bool(createPrFlagName) {
		config.PullRequest.Create = true
	}
	if v := ctx.String(prBodyFlagName); v != "" {
		config.PullRequest.BodyTemplate = v
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

func runUpdateCommand(*cli.Context) error {
	globalK := koanf.New(".")
	configDefaultName := "config_defaults.yml"

	if info, err := os.Stat(configDefaultName); err != nil || info.IsDir() {
		printer.WarnF("File %s does not exist, ignoring template defaults")
	} else {
		printer.DebugF("Loading config %s", configDefaultName)
		err = globalK.Load(file.Provider(configDefaultName), yaml.Parser())
		if err != nil {
			return nil
		}
	}
	var services []*repository.Service
	parser := rendering.NewParser(config.Template)

	logger := printer.PipelineLogger{Logger: printer.New().SetName("update").SetLevel(printer.DefaultLevel)}
	p := pipeline.NewPipelineWithLogger(logger).WithSteps(
		pipeline.NewStep("parse templates", parser.ParseTemplateDirAction()),
		pipeline.NewStep("parse managed repos config", parseServices(&services)),
		parallel.NewWorkerPoolStep("update repositories", 1, updateReposInParallel(&services, globalK, parser), errorHandler(&services)),
	)
	return p.Run().Err
}

func updateReposInParallel(services *[]*repository.Service, globalK *koanf.Koanf, parser *rendering.Parser) parallel.PipelineSupplier {
	return func(pipelines chan *pipeline.Pipeline) {
		defer close(pipelines)
		for _, r := range *services {
			p := createPipeline(r, globalK, parser)
			pipelines <- p
		}
	}
}

func createPipeline(r *repository.Service, globalK *koanf.Koanf, parser *rendering.Parser) *pipeline.Pipeline {
	log := printer.New().SetName(r.Config.Name).SetLevel(printer.DefaultLevel)

	sc := &cfg.SyncConfig{
		Git:         r.Config,
		PullRequest: config.PullRequest,
		Template: &cfg.TemplateConfig{
			RootDir: config.Template.RootDir,
		},
	}
	renderer := rendering.NewRenderer(sc, globalK, parser)
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
			predicate.ToStep("commit", r.Commit(), r.EnabledCommit()),
			predicate.ToStep("show diff", r.Diff(), r.EnabledCommit()),
			predicate.ToStep("push", r.PushToRemote(), r.EnabledPush()),
		).AsNestedStep("push changes"), r.Dirty()),
		predicate.WrapIn(pipeline.NewPipelineWithLogger(logger).WithSteps(
			pipeline.NewStep("render pull request template", renderer.RenderPrTemplate()),
			pipeline.NewStep("create or update pull request", r.CreateOrUpdatePr(config.PullRequest)),
		).AsNestedStep("pull request"), predicate.Bool(sc.PullRequest.Create)),
	)
	return p
}

func parseServices(services *[]*repository.Service) func() pipeline.Result {
	return func() pipeline.Result {
		*services = repository.NewServicesFromFile(config)
		return pipeline.Result{}
	}
}

func errorHandler(services *[]*repository.Service) parallel.ResultHandler {
	return func(results map[uint64]pipeline.Result) pipeline.Result {
		var err error
		for index, service := range *services {
			if result := results[uint64(index)]; result.Err != nil {
				err = multierror.Append(err, fmt.Errorf("%s: %w", service.Config.Name, result.Err))
			}
		}
		return pipeline.Result{Err: err}
	}
}
