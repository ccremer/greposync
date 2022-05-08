package flags

import (
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

const (
	// ProjectIncludeFlagName is the name on the CLI
	ProjectIncludeFlagName = "include"
	// ProjectExcludeFlagName is the name on the CLI
	ProjectExcludeFlagName = "exclude"
	// ProjectJobsFlagName is the name on the CLI
	ProjectJobsFlagName = "jobs"
)

var (
	JobsMinimumCount = 1
	JobsMaximumCount = 8
)

// EnvVarPrefix is the environment variable prefix key.
var EnvVarPrefix = "G_"

//// Log Flags

func NewLogLevelFlag(dst *int) *cli.IntFlag {
	return &cli.IntFlag{Name: "log.level", EnvVars: Prefixed("LOG_LEVEL"), Aliases: []string{"v"},
		Usage: "Log level that increases verbosity with greater numbers.",
		Value: 0, Destination: dst,
	}
}

func NewShowLogFlag(dst *bool) *altsrc.BoolFlag {
	return altsrc.NewBoolFlag(&cli.BoolFlag{Name: "log.showLog", EnvVars: Prefixed("SHOW_LOG"),
		Usage: "Shows the full log in real-time rather than keeping it hidden until an error occurred.",
		Value: false, Destination: dst,
	})
}

func NewShowDiffFlag(dst *bool) *altsrc.BoolFlag {
	return altsrc.NewBoolFlag(&cli.BoolFlag{Name: "log.showDiff", EnvVars: Prefixed("SHOW_DIFF"),
		Usage: "Show the Git Diff for each repository after committing. In --dry-run=offline mode the diff is showed for unstaged changes.",
		Value: false, Destination: dst,
	})
}

//// Common Flags

func NewJobsFlag(dst *int) *cli.IntFlag {
	return &cli.IntFlag{Name: ProjectJobsFlagName, EnvVars: Prefixed("JOBS"), Aliases: []string{"j"},
		Usage: "Jobs is the number of parallel jobs to run. 1 basically means that jobs are run in sequence.",
		Value: 1, Destination: dst,
	}
}

func NewSkipBrokenFlag(dst *bool) *cli.BoolFlag {
	return &cli.BoolFlag{Name: "skipBroken", EnvVars: Prefixed("SKIP_BROKEN"),
		Usage: "Skip abort if a repository update encounters an error",
		Value: false, Destination: dst,
	}
}

func NewIncludeFlag(dst *string) *cli.StringFlag {
	return &cli.StringFlag{Name: ProjectIncludeFlagName, EnvVars: Prefixed("INCLUDE"),
		Usage: "Includes only repositories in the update that match the given filter (regex). The full URL (including scheme) is matched.",
		Value: "", Destination: dst,
	}
}

func NewExcludeFlag(dst *string) *cli.StringFlag {
	return &cli.StringFlag{Name: ProjectExcludeFlagName, EnvVars: Prefixed("EXCLUDE"),
		Usage: "Excludes repositories from updating that match the given filter (regex). Repositories matching both include and exclude filter are still excluded.",
		Value: "", Destination: dst,
	}
}

func NewDryRunFlag(dst *string) *cli.StringFlag {
	return &cli.StringFlag{Name: "dry-run", EnvVars: Prefixed("DRYRUN"),
		Usage: "Select a dry run mode. Allowed values: offline (do not run any Git commands except initial clone), commit (commit, but don't push), push (push, but don't touch PRs)",
		Value: "", Destination: dst,
	}
}

//// Git Flags

func NewGitRootDirFlag(dst *string) *altsrc.PathFlag {
	return altsrc.NewPathFlag(&cli.PathFlag{Name: "git.root", EnvVars: Prefixed("GIT_ROOT_DIR"),
		Usage: "Local relative directory path where git clones repositories into.",
		Value: "repos", Destination: dst,
	})
}

func NewGitAmendFlag(dst *bool) *cli.BoolFlag {
	return &cli.BoolFlag{Name: "git.amend", EnvVars: Prefixed("GIT_AMEND"),
		Usage: "Amend previous commit. Requires --git.forcePush.",
		Value: false, Destination: dst,
	}
}

func NewGitForcePushFlag(dst *bool) *altsrc.BoolFlag {
	return altsrc.NewBoolFlag(&cli.BoolFlag{Name: "git.forcePush", EnvVars: Prefixed("GIT_FORCEPUSH"),
		Usage: "If push is enabled, push forcefully.",
		Value: false, Destination: dst,
	})
}

func NewGitDefaultNamespaceFlag(dst *string) *altsrc.StringFlag {
	return altsrc.NewStringFlag(&cli.StringFlag{Name: "git.defaultNamespace", EnvVars: Prefixed("GIT_DEFAULT_NS"),
		Usage: "The repository owner without the repository name. This is often a user or organization name in GitHub.com or GitLab.com.",
		Value: "github.com", Destination: dst,
	})
}

func NewGitCommitBranchFlag(dst *string) *altsrc.StringFlag {
	return altsrc.NewStringFlag(&cli.StringFlag{Name: "git.commitBranch", EnvVars: Prefixed("GIT_COMMIT_BRANCH"),
		Usage: "The branch name to create, switch to and commit locally.",
		Value: "greposync-update", Destination: dst,
	})
}

func NewGitCommitMessageFlag(dst *string) *altsrc.StringFlag {
	return altsrc.NewStringFlag(&cli.StringFlag{Name: "git.commitMessage", EnvVars: Prefixed("GIT_COMMIT_MSG"),
		Usage: "The commit message when committing an update.",
		Value: "Update from greposync", Destination: dst,
	})
}

func NewGitBaseURLFlag(dst *string) *altsrc.StringFlag {
	return altsrc.NewStringFlag(&cli.StringFlag{Name: "git.base", EnvVars: Prefixed("GIT_BASE"),
		Usage: "Git base URL.",
		Value: "git@github.com:", Destination: dst,
	})
}

//// PR Flags

func NewPRCreateFlag(dst *bool) *altsrc.BoolFlag {
	return altsrc.NewBoolFlag(&cli.BoolFlag{Name: "pr.create", EnvVars: Prefixed("PR_CREATE"),
		Usage: "Create a PullRequest on a supported git hoster after pushing to remote.",
		Value: false, Destination: dst,
	})
}

func NewPRBodyFlag(dst *string) *altsrc.StringFlag {
	return altsrc.NewStringFlag(&cli.StringFlag{Name: "pr.body", EnvVars: Prefixed("PR_BODY"),
		Usage: "Markdown-enabled body of the PullRequest. It will load from an existing file if this is a path. Content can be templated.",
		Value: "This Pull request updates this repository with changes from a greposync template repository.", Destination: dst,
	})
}

func NewPRSubjectFlag(dst *string) *altsrc.StringFlag {
	return altsrc.NewStringFlag(&cli.StringFlag{Name: "pr.subject", EnvVars: Prefixed("PR_SUBJECT"),
		Usage: "The Pull Request title.",
		Value: "Update from greposync", Destination: dst,
	})
}

func NewPRTargetBranchFlag(dst *string) *altsrc.StringFlag {
	return altsrc.NewStringFlag(&cli.StringFlag{Name: "pr.targetBranch", EnvVars: Prefixed("PR_TARGET_BRANCH"),
		Usage: "Remote branch name of the pull request. If left empty, it will target the default branch (usually 'master' or 'main').",
		Value: "", Destination: dst,
	})
}

func NewPRLabelsFlag(dst *cli.StringSlice) *altsrc.StringSliceFlag {
	return altsrc.NewStringSliceFlag(&cli.StringSliceFlag{Name: "pr.labels", EnvVars: Prefixed("PR_LABELS"),
		Usage: "Array of issue labels to apply when creating a pull request. Labels on existing pull requests are not updated. It is not validated whether the labels exist, the API may or may not create non-existing labels dynamically.",
		Value: &cli.StringSlice{}, Destination: dst,
	})
}

//// Other Flags

func NewTemplateRootDirFlag(dst *string) *altsrc.PathFlag {
	return altsrc.NewPathFlag(&cli.PathFlag{Name: "template.root", EnvVars: Prefixed("TEMPLATE_ROOT_DIR"),
		Usage: "The path relative to the current workdir where the template files are located.",
		Value: "template", Destination: dst,
	})
}
