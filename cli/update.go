package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/ccremer/git-repo-sync/cfg"
	"github.com/ccremer/git-repo-sync/printer"
	"github.com/ccremer/git-repo-sync/rendering"
	"github.com/ccremer/git-repo-sync/repository"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/urfave/cli/v2"
)

const (
	dryRunFlagName   = "dry-run"
	createPrFlagName = "pr"
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

	if ctx.IsSet(dryRunFlagName) {
		dryRunMode := ctx.String(dryRunFlagName)
		switch dryRunMode {
		case "offline":
			config.SkipReset = true
			config.SkipCommit = true
			config.SkipPush = true
			config.PullRequest.Create = false
		case "commit":
			config.SkipPush = true
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

	for _, repo := range services {
		repo.PrepareWorkspace()

		renderer := rendering.NewService(rendering.Config{
			TemplateDir:    config.TemplateRoot,
			TargetRootDir:  repo.Config.GitDir,
			ConfigDefaults: globalK,
			RepoName:       path.Base(repo.Config.GitDir),
		})

		renderer.ProcessTemplates()
		repo.MakeCommit()
		repo.ShowDiff()
		repo.PushToRemote()
		repo.CreatePR()
	}
	return nil
}
