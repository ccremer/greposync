package labels

import (
	"encoding/json"
	"regexp"

	"github.com/ccremer/greposync/application/clierror"
	"github.com/ccremer/greposync/application/flags"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/printer"
	"github.com/urfave/cli/v2"
)

func (c *Command) validateCommand(ctx *cli.Context) error {
	if err := cfg.ParseConfig(c.cfg.Project.MainConfigFileName, c.cfg, ctx); err != nil {
		return clierror.AsUsageError(err)
	}

	if err := flags.ValidateGlobalFlags(ctx, c.cfg); err != nil {
		return clierror.AsUsageError(err)
	}

	if _, err := regexp.Compile(c.cfg.Project.Include); err != nil {
		return clierror.AsUsageErrorf("invalid flag --%s: %v", flags.ProjectIncludeFlagName, err)
	}
	if _, err := regexp.Compile(c.cfg.Project.Exclude); err != nil {
		return clierror.AsUsageErrorf("invalid flag --%s: %v", flags.ProjectExcludeFlagName, err)
	}

	if jobs := c.cfg.Project.Jobs; jobs > flags.JobsMaximumCount || jobs < flags.JobsMinimumCount {
		return clierror.AsFlagUsageErrorf(flags.ProjectJobsFlagName, "value is not between %d and %d", flags.JobsMinimumCount, flags.JobsMaximumCount)
	}

	_, err := cfg.RepositoryLabelSetConverter{}.ConvertToEntity(c.cfg.RepositoryLabels.Values())
	if err != nil {
		return clierror.AsUsageErrorf("invalid label configuration in '%s': %w", "repositoryLabels", err)
	}

	c.cfg.Sanitize()
	j, _ := json.Marshal(c.cfg)
	printer.DebugF("Using config: %s", j)
	return nil
}
