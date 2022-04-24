package test

import (
	"encoding/json"
	"regexp"

	"github.com/ccremer/greposync/application/clierror"
	"github.com/ccremer/greposync/application/flags"
	"github.com/ccremer/greposync/cfg"
	"github.com/urfave/cli/v2"
)

func (c *Command) validateTestCommand(ctx *cli.Context) error {
	if err := cfg.ParseConfig(c.cfg.Project.MainConfigFileName, c.cfg, ctx); err != nil {
		return clierror.AsUsageError(err)
	}

	if _, err := regexp.Compile(c.cfg.Project.Include); err != nil {
		return clierror.AsFlagUsageError(flags.ProjectIncludeFlagName, err)
	}
	if _, err := regexp.Compile(c.cfg.Project.Exclude); err != nil {
		return clierror.AsFlagUsageError(flags.ProjectExcludeFlagName, err)
	}

	if jobs := c.cfg.Project.Jobs; jobs > flags.JobsMaximumCount || jobs < flags.JobsMinimumCount {
		return clierror.AsFlagUsageErrorf(flags.ProjectJobsFlagName, "value is not between %d and %d", flags.JobsMinimumCount, flags.JobsMaximumCount)
	}

	c.logFactory.SetLogLevel(c.cfg.Log.Level)
	j, _ := json.Marshal(c.cfg)
	c.logFactory.NewGenericLogger("").V(1).Info("Using config", "config", string(j))
	return nil
}
