package github

import (
	"github.com/ccremer/greposync/cfg"
	"github.com/google/go-github/v37/github"
)

type (
	LabelConfig struct {
		Labels []cfg.RepositoryLabel
	}
)

func (p *Provider) UpdateLabels(c LabelConfig) error {
	if len(c.Labels) == 0 {
		p.log.InfoF("No labels defined in config, nothing to do.")
		return nil
	}
	nextPage := 1
	lastPage := 1
	var allLabels []*github.Label
	for repeat := true; repeat; repeat = nextPage < lastPage {
		labels, resp, err := p.client.Issues.ListLabels(p.ctx, p.cfg.RepoOwner, p.cfg.Repo, &github.ListOptions{
			Page:    1,
			PerPage: 100,
		})
		if err != nil {
			return err
		}
		p.setRemainingApiCalls(resp)
		allLabels = append(allLabels, labels...)
		lastPage = resp.LastPage
		nextPage = resp.NextPage
	}

	for _, label := range c.Labels {
		ghLabel := p.findGhLabel(allLabels, label)
		err := p.upsertLabel(ghLabel, label)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) upsertLabel(ghLabel *github.Label, label cfg.RepositoryLabel) error {
	var resp *github.Response
	var err error
	if ghLabel != nil && p.hasLabelChanged(ghLabel, label) {
		_, resp, err = p.client.Issues.EditLabel(p.ctx, p.cfg.RepoOwner, p.cfg.Repo, label.Name, ghLabel)
		p.setRemainingApiCalls(resp)
		if err != nil {
			return err
		}
		p.log.InfoF("Label '%s' updated", label.Name)
	} else if ghLabel == nil {
		_, resp, err = p.client.Issues.CreateLabel(p.ctx, p.cfg.RepoOwner, p.cfg.Repo, &github.Label{
			Name:        &label.Name,
			Color:       &label.Color,
			Description: &label.Description,
		})
		p.setRemainingApiCalls(resp)
		if err != nil {
			return err
		}
		p.log.InfoF("Label '%s' created", label.Name)
	} else if label.Delete {
		resp, err = p.client.Issues.DeleteLabel(p.ctx, p.cfg.RepoOwner, p.cfg.Repo, label.Name)
		p.setRemainingApiCalls(resp)
		if err != nil {
			return err
		}
		p.log.InfoF("Label '%s' deleted", label.Name)
	} else {
		p.log.InfoF("Label '%s' unchanged", label.Name)
	}
	return nil
}

func (p *Provider) findGhLabel(ghLabels []*github.Label, toFind cfg.RepositoryLabel) *github.Label {
	for _, label := range ghLabels {
		if label.GetName() == toFind.Name {
			return label
		}
	}
	return nil
}

func (p *Provider) hasLabelChanged(ghLabel *github.Label, repoLabel cfg.RepositoryLabel) bool {
	changed := false
	if ghLabel.GetDescription() != repoLabel.Description {
		ghLabel.Description = &repoLabel.Description
		changed = true
	}
	if ghLabel.GetColor() != repoLabel.Color {
		ghLabel.Color = &repoLabel.Color
		changed = true
	}
	return changed
}
