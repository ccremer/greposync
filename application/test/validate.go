package test

import (
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

	c.appService.console.SetTitle("RUNNING TESTS...")
	c.appService.console.SetCommandName("Test")
	c.appService.templateStore.SkipRemovingFileExtension = true
	c.appService.repoStore.ParentDir = "tests"
	c.appService.repoStore.TestOutputRootDir = ".tests"
	c.appService.repoStore.DefaultNamespace = "local"
	c.appService.engine.RootDir = c.appService.templateStore.RootDir
	c.logFactory.SetLogLevel(c.cfg.Log.Level)
	c.logFactory.NewGenericLogger("").V(1).Info("Using config", "config", flags.CollectFlagValues(ctx))
	return nil
}
