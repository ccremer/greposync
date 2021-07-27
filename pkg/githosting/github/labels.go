package github

import (
	"time"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
	"github.com/google/go-github/v37/github"
)

func (p *Facade) CreateOrUpdateLabelsForRepo(url *core.GitUrl, labels []core.GitRepositoryLabel) error {
	converted := LabelConverter{}.convertFromEntity(labels)
	p.log.SetName(url.GetRepositoryName())

	allLabels, err := p.fetchAllLabels(url)
	if err != nil {
		return err
	}

	for _, label := range converted {
		ghLabel := p.findMatchingGhLabel(allLabels, label)
		if ghLabel == nil {
			createErr := p.createLabel(url, label)
			if createErr == nil {
				p.log.InfoF("Label '%s' created", label.Name)
			}
			return createErr
		}
		if !p.hasLabelChanged(ghLabel, label) {
			p.log.InfoF("Label '%s' unchanged", label.Name)
			continue
		}
		err = p.updateLabel(url, ghLabel, label)
		if err != nil {
			return err
		}
		p.log.InfoF("Label '%s' updated", label.Name)
	}

	return nil
}

func (p *Facade) DeleteLabelsForRepo(url *core.GitUrl, labels []core.GitRepositoryLabel) error {
	p.log.SetName(url.GetRepositoryName())
	converted := LabelConverter{}.convertFromEntity(labels)
	for _, label := range converted {
		err := p.deleteLabel(url, label)
		if err != nil {
			return err
		}
		p.log.InfoF("Label '%s' deleted", label.Name)
		/*
			https://docs.github.com/en/rest/guides/best-practices-for-integrators#dealing-with-abuse-rate-limits
			"If you're making a large number of POST, PATCH, PUT, or DELETE requests for a single user or client ID, wait at least one second between each request."
		*/
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (p *Facade) createLabel(url *core.GitUrl, label *cfg.RepositoryLabel) error {
	_, _, err := p.client.Issues.CreateLabel(p.ctx, url.GetNamespace(), url.GetRepositoryName(), &github.Label{
		Name:        &label.Name,
		Color:       &label.Color,
		Description: &label.Description,
	})
	return err
}

func (p *Facade) updateLabel(url *core.GitUrl, ghLabel *github.Label, label *cfg.RepositoryLabel) error {
	// TODO: Without a new_name property we cannot rename a label yet.
	ghLabel.Description = &label.Description
	ghLabel.Color = &label.Color
	_, _, err := p.client.Issues.EditLabel(p.ctx, url.GetNamespace(), url.GetRepositoryName(), label.Name, ghLabel)
	return err
}

func (p *Facade) deleteLabel(url *core.GitUrl, label *cfg.RepositoryLabel) error {
	_, err := p.client.Issues.DeleteLabel(p.ctx, url.GetNamespace(), url.GetRepositoryName(), label.Name)
	return err
}

func (p *Facade) fetchAllLabels(url *core.GitUrl) ([]*github.Label, error) {
	nextPage := 1
	lastPage := 1
	var allLabels []*github.Label
	for repeat := true; repeat; repeat = nextPage < lastPage {
		labels, resp, err := p.client.Issues.ListLabels(p.ctx, url.GetNamespace(), url.GetRepositoryName(), &github.ListOptions{
			Page:    1,
			PerPage: 100,
		})
		if err != nil {
			return nil, err
		}
		allLabels = append(allLabels, labels...)
		lastPage = resp.LastPage
		nextPage = resp.NextPage
	}
	return allLabels, nil
}

func (p *Facade) findMatchingGhLabel(ghLabels []*github.Label, toFind *cfg.RepositoryLabel) *github.Label {
	for _, label := range ghLabels {
		if label.GetName() == toFind.Name {
			return label
		}
	}
	return nil
}

func (p *Facade) hasLabelChanged(ghLabel *github.Label, repoLabel *cfg.RepositoryLabel) bool {
	return ghLabel.GetDescription() != repoLabel.Description || ghLabel.GetColor() != repoLabel.Color
}
