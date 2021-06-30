package cfg

import (
	"net/url"
	"path"
)

type (
	// Configuration holds a strongly-typed tree of the main configuration
	Configuration struct {
		ProjectRoot string
		Log         LogConfig
		PullRequest PullRequestConfig
		Template    TemplateConfig
		Git         GitConfig
	}
	// LogConfig configures the logging options
	LogConfig struct {
		Level string
	}
	// PullRequestConfig configures the pull request feature
	PullRequestConfig struct {
		Create       bool
		TargetBranch string
		Labels       []string
		BodyTemplate string
		Subject      string
	}

	// SyncConfig configures a single repository sync
	SyncConfig struct {
		PullRequest *PullRequestConfig
		Git         *GitConfig
		Template    *TemplateConfig
	}
	GitConfig struct {
		Url           *url.URL `json:"-"`
		Dir           string   `json:"-"`
		SkipReset     bool
		SkipCommit    bool
		SkipPush      bool
		ForcePush     bool
		CreatePR      bool
		Amend         bool `json:"-"`
		CommitMessage string
		CommitBranch  string
		DefaultBranch string
		Name          string
		Namespace     string
	}
	TemplateConfig struct {
		RootDir string
	}
)

// NewDefaultConfig retrieves the hardcoded configs with sane defaults
func NewDefaultConfig() *Configuration {
	return &Configuration{
		ProjectRoot: "repos",
		Log: LogConfig{
			Level: "info",
		},
		Git: GitConfig{
			CommitMessage: "Update from git-repo-sync",
		},
		PullRequest: PullRequestConfig{
			BodyTemplate: `This Pull request updates this repository with changes from a git-repo-sync template repository.`,
			Subject:      "Update from git-repo-sync",
		},
		Template: TemplateConfig{
			RootDir: "template",
		},
	}
}

func (c GitConfig) GetName() string {
	return path.Base(c.Dir)
}
