package cli

import (
	"encoding/json"
	"fmt"
	"os"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/printer"
	"github.com/ccremer/greposync/rendering"
	"github.com/ccremer/greposync/repository"
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
	services := repository.NewServicesFromFile(config)
	parser := rendering.NewParser(config.Template)

	if err := parser.ParseTemplateDir(); err != nil {
		return err
	}

	for _, r := range services {
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
				pipeline.NewStepWithPredicate("clone repository", r.CloneGitRepository(), pipeline.Bool(!gitDirExists)),
				pipeline.NewStep("determine default branch", r.GetDefaultBranch()),
				pipeline.NewStepWithPredicate("fetch", r.Fetch(), r.EnabledReset()),
				pipeline.NewStepWithPredicate("reset repository", r.ResetRepository(), r.EnabledReset()),
				pipeline.NewStep("checkout branch", r.CheckoutBranch()),
				pipeline.NewStepWithPredicate("pull", r.Pull(), r.EnabledReset()),
			).AsNestedStep("prepare workspace", nil),
			pipeline.NewStep("render templates", renderer.RenderTemplateDir()),
			pipeline.NewPipelineWithLogger(logger).WithSteps(
				pipeline.NewStepWithPredicate("add", r.Add(), r.EnabledCommit()),
				pipeline.NewStepWithPredicate("commit", r.Commit(), r.EnabledCommit()),
				pipeline.NewStepWithPredicate("show diff", r.Diff(), r.EnabledCommit()),
				pipeline.NewStepWithPredicate("push", r.PushToRemote(), r.EnabledPush()),
			).AsNestedStep("push changes", r.Dirty()),
			pipeline.NewPipelineWithLogger(logger).WithSteps(
				pipeline.NewStep("render pull request template", renderer.RenderPrTemplate()),
				pipeline.NewStep("create or update pull request", r.CreateOrUpdatePr(config.PullRequest)),
			).AsNestedStep("pull request", pipeline.Bool(sc.PullRequest.Create)),
		)
		result := p.Run()
		if !result.IsSuccessful() {
			return result.Err
		}
	}
	return nil
}
