package cfg

type (
	// Configuration holds a strongly-typed tree of the main configuration
	Configuration struct {
		Project          *ProjectConfig     `json:"project" koanf:"project"`
		Log              *LogConfig         `json:"log" koanf:"log"`
		PullRequest      *PullRequestConfig `json:"pr" koanf:"pr"`
		Template         *TemplateConfig    `json:"template" koanf:"template"`
		Git              *GitConfig         `json:"git" koanf:"git"`
		RepositoryLabels RepositoryLabelMap `json:"repositoryLabels" koanf:"repositoryLabels"`
	}
	// ProjectConfig configures the main config settings
	ProjectConfig struct {
		MainConfigFileName    string `json:"-"`
		ConfigDefaultFileName string `json:"-"`

		// RootDir is the local directory where the Git repositories are cloned into.
		RootDir string `json:"rootDir" koanf:"rootDir"`
		// Jobs is the number of parallel jobs to run.
		// Requires a minimum of 1, supports a maximum of 8.
		// 1 basically means that jobs are run in sequence.
		// If this number is 2 or greater, then the logs are buffered and only displayed in case of errors.
		Jobs int `json:"jobs" koanf:"jobs"`
		// Include is a regex filter that includes repositories only when they match.
		// The filter is applied to the whole URL.
		// This option is not configurable in `greposync.yml`.
		Include string `json:"-" koanf:"include"`
		// Exclude is similar to Include, only that matching repository URLs are skipped.
		// This option is not configurable in `greposync.yml`.
		Exclude string `json:"-" koanf:"exclude"`
	}

	// LogConfig configures the logging options
	LogConfig struct {
		Level    string `json:"level" koanf:"level"`
		ShowDiff bool   `json:"showDiff" koanf:"showDiff"`
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
	// RepositoryLabel is a struct describing a Label on a Git hosting service like GitHub.
	RepositoryLabel struct {
		// Name is the label name.
		Name string `json:"name" koanf:"name"`
		// Description is a short description of the label.
		Description string `json:"description" koanf:"description"`
		// Color is the hexadecimal color code for the label, without the leading #.
		Color string `json:"color" koanf:"color"`
		// Delete will remove this label.
		Delete bool `json:"delete" koanf:"delete"`
	}
	RepositoryLabelMap map[string]RepositoryLabel

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
		SkipReset  bool `json:"skipReset"`
		SkipCommit bool `json:"skipCommit"`
		SkipPush   bool `json:"skipPush"`
		ForcePush  bool `json:"forcePush"`
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
			RootDir:               "repos",
			Jobs:                  1,
			MainConfigFileName:    "greposync.yml",
			ConfigDefaultFileName: "config_defaults.yml",
		},
		Log: &LogConfig{
			Level: "error",
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

func (s RepositoryLabelMap) Values() []RepositoryLabel {
	list := make([]RepositoryLabel, 0)
	for _, label := range s {
		list = append(list, label)
	}
	return list
}

func (s RepositoryLabelMap) SelectDeletions() []RepositoryLabel {
	list := make([]RepositoryLabel, 0)
	for _, label := range s {
		if label.Delete {
			list = append(list, label)
		}
	}
	return list
}

func (s RepositoryLabelMap) SelectModifications() []RepositoryLabel {
	list := make([]RepositoryLabel, 0)
	for _, label := range s {
		if !label.Delete {
			list = append(list, label)
		}
	}
	return list
}
