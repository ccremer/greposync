package cfg

type (
	// Configuration holds a strongly-typed tree of the configuration
	Configuration struct {
		Log         LogConfig
		Namespace   string
		Message     string
		PullRequest PullRequestConfig
		SkipCommit  bool
		SkipPush    bool
		ProjectRoot string
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
	}
)

// NewDefaultConfig retrieves the hardcoded configs with sane defaults
func NewDefaultConfig() *Configuration {
	return &Configuration{
		Log: LogConfig{
			Level: "info",
		},
		Message:     "Update from git-repo-sync",
		ProjectRoot: "repos",
	}
}
