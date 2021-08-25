package repository

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/infrastructure/githosting"
	"github.com/ccremer/greposync/printer"
)

type (
	// Repository is the implementation for core.GitRepository.
	Repository struct {
		GitConfig  *cfg.GitConfig
		PrConfig   *cfg.PullRequestConfig
		coreConfig core.GitRepositoryProperties
		labels []core.Label
		remote githosting.Remote
		pr     core.PullRequest
		log        printer.Printer
	}
)

func NewGitRepository(cfg *cfg.GitConfig, prConfig *cfg.PullRequestConfig, labels []core.Label) *Repository {
	return &Repository{
		GitConfig: cfg,
		coreConfig: core.GitRepositoryProperties{
			URL:     core.FromURL(cfg.Url),
			RootDir: cfg.Dir,
		},
		PrConfig: prConfig,
		labels:   labels,
		log:      printer.New().SetName(cfg.Name),
	}
}

// GetLabels implements core.GitRepository.
func (s *Repository) GetLabels() []core.Label {
	return s.labels
}
