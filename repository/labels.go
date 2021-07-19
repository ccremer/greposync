package repository

import (
	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/repository/github"
)

// CreateOrUpdateLabels creates or updates the labels configured.
func (s *Service) CreateOrUpdateLabels(config map[string]cfg.RepositoryLabel) pipeline.ActionFunc {
	return func() pipeline.Result {
		lc := github.LabelConfig{Labels: toArray(config)}
		return pipeline.Result{Err: s.provider.UpdateLabels(lc)}
	}
}

func toArray(labels map[string]cfg.RepositoryLabel) []cfg.RepositoryLabel {
	var list []cfg.RepositoryLabel
	for _, v := range labels {
		list = append(list, v)
	}
	return list
}
