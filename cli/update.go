package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"path"

	"github.com/ccremer/git-repo-sync/cfg"
	"github.com/ccremer/git-repo-sync/printer"
	"github.com/ccremer/git-repo-sync/rendering"
	"github.com/ccremer/git-repo-sync/repository"
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

	services := repository.NewServicesFromFile("managed_repos.yml", config.ProjectRoot, config.Namespace)

	for _, repoService := range services {
		repoService.PrepareWorkspace()

		rendering.LoadConfigFile("config_defaults.yml")
		syncFile := path.Join(repoService.Config.GitDir, ".sync.yml")
		rendering.LoadConfigFile(syncFile)

		data := map[string]interface{}{
			"Values": rendering.Unmarshal("README.md/test"),
		}

		err := rendering.RenderTemplate(repoService.Config.GitDir, data)
		if err != nil {
			log.Fatal(err)
		}

		repoService.MakeCommit()
		repoService.ShowDiff()
		repoService.PushToRemote()
		repoService.CreatePR()
	}
	return nil
}
