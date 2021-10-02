package github

import (
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/google/go-github/v39/github"
)

// GitHubInstrumentation is responsible for logging interactions with GitHub API.
type GitHubInstrumentation struct {
	factory logging.LoggerFactory
}

// NewGitHubInstrumentation returns a new instance.
func NewGitHubInstrumentation(factory logging.LoggerFactory) *GitHubInstrumentation {
	return &GitHubInstrumentation{
		factory: factory,
	}
}

func (i *GitHubInstrumentation) fetchedAllLabels(repository *domain.GitRepository, labels []*github.Label) {
	log := i.factory.NewRepositoryLogger(repository).V(logging.LevelDebug)
	if log.Enabled() {
		labelArr := make([]string, len(labels))
		for i, label := range labels {
			labelArr[i] = label.GetName()
		}
		log.Info("Fetched labels", "labels", labelArr)
	}
}

func (i *GitHubInstrumentation) createdLabel(repository *domain.GitRepository, label domain.Label, err error) error {
	if err == nil {
		i.factory.NewRepositoryLogger(repository).V(logging.LevelDebug).Info("Created label", "label", label.Name)
	}
	return err
}

func (i *GitHubInstrumentation) updatedLabel(repository *domain.GitRepository, label domain.Label, err error) error {
	if err == nil {
		i.factory.NewRepositoryLogger(repository).V(logging.LevelDebug).Info("Updated label", "label", label.Name)
	}
	return err
}

func (i *GitHubInstrumentation) deletedLabel(repository *domain.GitRepository, label *github.Label, err error) error {
	if err == nil {
		i.factory.NewRepositoryLogger(repository).V(logging.LevelDebug).Info("Deleted label", "label", label.Name)
	}
	return err
}

func (i *GitHubInstrumentation) prCreated(repository *domain.GitRepository, htmlUrl string) {
	i.factory.NewRepositoryLogger(repository).V(logging.LevelInfo).Info("PR created", "url", htmlUrl)
}

func (i *GitHubInstrumentation) prNotCreatedBecauseNoCommits(repository *domain.GitRepository, pr *domain.PullRequest) error {
	i.factory.NewRepositoryLogger(repository).V(logging.LevelInfo).Info("No pull request created as there are no commits between branches", "base", pr.BaseBranch, "head", pr.CommitBranch)
	return nil
}

func (i *GitHubInstrumentation) noPrFound(repository *domain.GitRepository) error {
	i.factory.NewRepositoryLogger(repository).V(logging.LevelDebug).Info("No PR found")
	return nil
}

func (i *GitHubInstrumentation) prFound(repository *domain.GitRepository, pr *github.PullRequest) error {
	i.factory.NewRepositoryLogger(repository).V(logging.LevelDebug).Info("Existing PR found", "url", pr.GetHTMLURL())
	return nil
}

func (i *GitHubInstrumentation) prUpdated(repository *domain.GitRepository, pr *github.PullRequest, err error) error {
	if err == nil {
		i.factory.NewRepositoryLogger(repository).V(logging.LevelInfo).Info("Updated pull request", "url", pr.GetHTMLURL(), "title", pr.GetTitle())
	}
	return err
}

func (i *GitHubInstrumentation) prLabelsUpdated(repository *domain.GitRepository, pr *domain.PullRequest, err error) error {
	if err == nil {
		i.factory.NewRepositoryLogger(repository).V(logging.LevelDebug).Info("Updated pull request labels", "labels", pr.GetLabels().String())
	}
	return err
}

func (i *GitHubInstrumentation) prIsUpToDate(repository *domain.GitRepository, cached *github.PullRequest) error {
	i.factory.NewRepositoryLogger(repository).V(logging.LevelInfo).Info("Pull request is up-to-date", "url", cached.GetHTMLURL())
	return nil
}
