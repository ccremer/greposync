package update

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/ccremer/greposync/cfg"
	app "github.com/ccremer/greposync/cli"
	"github.com/ccremer/greposync/printer"
	"github.com/urfave/cli/v2"
)

func (c *Command) validateUpdateCommand(ctx *cli.Context) error {
	if err := cfg.ParseConfig(app.GrepoSyncFileName, c.cfg, ctx); err != nil {
		return err
	}

	if err := app.ValidateGlobalFlags(ctx); err != nil {
		return err
	}

	if _, err := regexp.Compile(c.cfg.Project.Include); err != nil {
		return fmt.Errorf("invalid flag --%s: %v", app.ProjectIncludeFlagName, err)
	}
	if _, err := regexp.Compile(c.cfg.Project.Exclude); err != nil {
		return fmt.Errorf("invalid flag --%s: %v", app.ProjectExcludeFlagName, err)
	}

	if jobs := c.cfg.Project.Jobs; jobs > app.JobsMaximumCount || jobs < app.JobsMinimumCount {
		return fmt.Errorf("--%s is required to be between %d and %d", app.ProjectJobsFlagName, app.JobsMinimumCount, app.JobsMaximumCount)
	}

	if ctx.IsSet(dryRunFlagName) {
		dryRunMode := ctx.String(dryRunFlagName)
		switch dryRunMode {
		case "offline":
			c.cfg.Git.SkipReset = true
			c.cfg.Git.SkipCommit = true
			c.cfg.Git.SkipPush = true
			c.cfg.PullRequest.Create = false
		case "commit":
			c.cfg.Git.SkipPush = true
			c.cfg.PullRequest.Create = false
		case "push":
			c.cfg.PullRequest.Create = false
		default:
			return fmt.Errorf("invalid flag value of %s: %s", dryRunFlagName, dryRunMode)
		}
	}

	c.cfg.Sanitize()
	j, _ := json.Marshal(c.cfg)
	printer.DebugF("Using config: %s", j)
	return nil
}
