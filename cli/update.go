package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ccremer/git-repo-sync/cfg"
	"github.com/ccremer/git-repo-sync/printer"
	"github.com/ccremer/git-repo-sync/rendering"
	"github.com/ccremer/git-repo-sync/repository"
	pipeline "github.com/ccremer/go-command-pipeline"
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

	for _, r := range services {
		log := printer.New().SetName(r.Config.Name).SetLevel(printer.LevelDebug)

		sc := &cfg.SyncConfig{
			Git:         r.Config,
			PullRequest: config.PullRequest,
			Template: &cfg.TemplateConfig{
				RootDir: config.Template.RootDir,
			},
		}
		gitDirExists := r.DirExists(r.Config.Dir)
		p := pipeline.NewPipelineWithLogger(printer.PipelineLogger{Logger: log})
		p.WithSteps(
			pipeline.NewStepWithPredicate("clone repository", r.CloneGitRepositoryAction(), pipeline.Bool(!gitDirExists)),
			pipeline.NewStepWithPredicate("reset repository", r.ResetRepository(), pipeline.Not(r.SkipReset())),
			pipeline.NewStep("checkout branch", r.CheckoutBranch()),
			pipeline.NewStepWithPredicate("pull", r.Pull(), pipeline.Not(r.SkipReset())),
		)
		result := p.Run()
		log.CheckIfError(result.Err)

		renderer := rendering.NewRenderer(sc, globalK)
		renderer.ProcessTemplates()

		p = pipeline.NewPipeline()
		p.WithSteps(
			pipeline.NewStepWithPredicate("commit", r.MakeCommit(), pipeline.Not(r.SkipCommit())),
			pipeline.NewStepWithPredicate("show diff", r.ShowDiff(), pipeline.Not(r.SkipCommit())),
			pipeline.NewStepWithPredicate("push", r.PushToRemote(), pipeline.Not(r.SkipPush())),
		)
		result = p.Run()
		log.CheckIfError(result.Err)

		if config.PullRequest.Create {
			template := config.PullRequest.BodyTemplate
			if template == "" {
				log.InfoF("No PullRequest template defined")
				config.PullRequest.BodyTemplate = config.Git.CommitMessage
			}

			data := rendering.Values{"Metadata": renderer.ConstructMetadata()}
			if renderer.FileExists(template) {
				config.PullRequest.BodyTemplate = renderer.RenderTemplateFile(data, template)
			} else {
				config.PullRequest.BodyTemplate = renderer.RenderString(data, template)
			}
			r.CreateOrUpdatePR(config.PullRequest)
		}
	}
	return nil
}
