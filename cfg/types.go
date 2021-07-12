package cfg

import (
	"net/url"
)

type (
	// Configuration holds a strongly-typed tree of the main configuration
	Configuration struct {
		Project     *ProjectConfig     `json:"project" koanf:"project"`
		Log         *LogConfig         `json:"log" koanf:"log"`
		PullRequest *PullRequestConfig `json:"pr" koanf:"pr"`
		Template    *TemplateConfig    `json:"template" koanf:"template"`
		Git         *GitConfig         `json:"git" koanf:"git"`
	}
	// ProjectConfig configures the project
	ProjectConfig struct {
		// RootDir is the local directory where the Git repositories are cloned into.
		RootDir string `json:"rootDir" koanf:"rootDir"`
		// Jobs is the number of parallel jobs to run.
		// Requires a minimum of 1, supports a maximum of 8.
		// 1 basically means that jobs are run in sequence.
		// If this number is 2 or greater, then the logs are buffered and only displayed in case of errors.
		Jobs int `json:"jobs" koanf:"jobs"`
	}

	// LogConfig configures the logging options
	LogConfig struct {
		Level string `json:"level" koanf:"level"`
	}
	// PullRequestConfig configures the pull request feature
	PullRequestConfig struct {
		Create bool `json:"create" koanf:"create"`
		// TargetBranch is the target remote branch of the pull request.
		// If left empty, it will target the default branch.
		TargetBranch string `json:"targetBranch" koanf:"targetBranch"`
		// Labels is an array of issue labels to apply when creating a pull request.
		// Labels on existing pull requests are not updated.
		// It is not validated whether the labels exist, the API may or may not create non-existing labels dynamically.
		Labels []string `json:"labels" koanf:"labels"`
		// BodyTemplate is the description used in pull requests.
		// Supports Go template with the `.Metadata` key.
		// If this string is a relative path to an existing file in the greposync directory, the file is parsed as a Go template.
		// If empty, the CommitMessage is used.
		BodyTemplate string `json:"bodyTemplate" koanf:"bodyTemplate"`
		// Subject is the Pull Request title.
		Subject string `json:"subject" koanf:"subject"`
	}

	// SyncConfig configures a single repository sync
	SyncConfig struct {
		PullRequest *PullRequestConfig `json:"pullRequest"`
		Git         *GitConfig         `json:"git"`
		Template    *TemplateConfig    `json:"template"`
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
		SkipReset  bool   `json:"skipReset"`
		SkipCommit bool   `json:"skipCommit"`
		SkipPush   bool   `json:"skipPush"`
		ForcePush  bool   `json:"forcePush"`
		// Amend will amend the last commit.
		// This option is not configurable in `greposync.yml`.
		// Configurable only via environment variables or CLI flag.
		Amend bool `json:"-" koanf:"amend"`
		// CommitMessage is the string that is passed to `git commit`.
		// It can contain newlines, for example to pass a long description.
		CommitMessage string `json:"commitMessage" koanf:"commitMessage"`
		CommitBranch  string `json:"commitBranch" koanf:"commitBranch"`
		// DefaultBranch is the name of the default branch in origin.
		DefaultBranch string `json:"defaultBranch"`
		// Name is the git repository name without .git extension.
		Name string `json:"name"`
		// Namespace is the repository owner without the repository name.
		// This is often a user or organization name in GitHub.com or GitLab.com.
		Namespace string `json:"namespace"`
	}
	// TemplateConfig configures template settings
	TemplateConfig struct {
		// RootDir is the path relative to the current workdir where the template files are located.
		RootDir string `json:"rootDir" koanf:"rootDir"`
	}
)

// NewDefaultConfig retrieves the hardcoded configs with sane defaults
func NewDefaultConfig() *Configuration {
	return &Configuration{
		Project: &ProjectConfig{
			RootDir: "repos",
			Jobs:    1,
		},
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
