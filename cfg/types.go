package cfg

import (
	"net/url"
)

type (
	// Configuration holds a strongly-typed tree of the main configuration
	Configuration struct {
		ProjectRoot string
		Log         *LogConfig
		PullRequest *PullRequestConfig
		Template    *TemplateConfig
		Git         *GitConfig
	}
	// LogConfig configures the logging options
	LogConfig struct {
		Level string
	}
	// PullRequestConfig configures the pull request feature
	PullRequestConfig struct {
		Create bool
		// TargetBranch is the target remote branch of the pull request.
		// If left empty, it will target the default branch.
		TargetBranch string
		// Labels is an array of issue labels to apply when creating a pull request.
		// Labels on existing pull requests are not updated.
		// It is not validated whether the labels exist, the API may or may not create non-existing labels dynamically.
		Labels []string
		// BodyTemplate is the description used in pull requests.
		// Supports Go template with the `.Metadata` key.
		// If this string is a relative path to an existing file in the greposync directory, the file is parsed as a Go template.
		// If empty, the CommitMessage is used.
		BodyTemplate string
		// Subject is the Pull Request title.
		Subject string
	}

	// SyncConfig configures a single repository sync
	SyncConfig struct {
		PullRequest *PullRequestConfig
		Git         *GitConfig
		Template    *TemplateConfig
	}
	// GitConfig configures a git repository.
	// This structure is used to configuring the sync behaviour
	// It is also passed to templates with filled-in information
	GitConfig struct {
		// Url is the full Git URL to the remote repository.
		// This option is not configurable in `greposync.yml`.
		// In templates, the URL is looking like ``
		Url *url.URL `json:"-"`
		// Dir is the relative path to current working directory where the repository is cloned locally.
		// This option is not configurable in `greposync.yml`.
		Dir        string `json:"-"`
		SkipReset  bool
		SkipCommit bool
		SkipPush   bool
		ForcePush  bool
		// Amend will amend the last commit.
		// This option is not configurable in `greposync.yml`.
		// Configurable only via environment variables or CLI flag.
		Amend bool `json:"-"`
		// CommitMessage is the string that is passed to `git commit`.
		// It can contain newlines, for example to pass a long description.
		CommitMessage string
		CommitBranch  string
		// DefaultBranch is the name of the default branch in origin.
		DefaultBranch string
		// Name is the git repository name without .git extension.
		Name string
		// Namespace is the repository owner without the repository name.
		// This is often a user or organization name in GitHub.com or GitLab.com.
		Namespace string
	}
	// TemplateConfig configures template settings
	TemplateConfig struct {
		// RootDir is the path relative to the current workdir where the template files are located.
		RootDir string
	}
)

// NewDefaultConfig retrieves the hardcoded configs with sane defaults
func NewDefaultConfig() *Configuration {
	return &Configuration{
		ProjectRoot: "repos",
		Log: &LogConfig{
			Level: "info",
		},
		Git: &GitConfig{
			CommitMessage: "Update from greposync",
		},
		PullRequest: &PullRequestConfig{
			BodyTemplate: `This Pull request updates this repository with changes from a greposync template repository.`,
			Subject:      "Update from greposync",
		},
		Template: &TemplateConfig{
			RootDir: "template",
		},
	}
}
