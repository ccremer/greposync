package github

import (
	"time"

	"github.com/ccremer/greposync/cfg"
	"github.com/ccremer/greposync/core"
	"github.com/google/go-github/v37/github"
)

// CreateOrUpdateLabelsForRepo implements core.GitHostingFacade.
func (p *Facade) CreateOrUpdateLabelsForRepo(url *core.GitURL, labels []core.Label) error {
	converted := LabelConverter{}.ConvertFromEntity(labels)
	p.log.SetName(url.GetRepositoryName())

	allLabels, err := p.fetchAllLabels(url)
	if err != nil {
		return err
	}

	for _, label := range converted {
		ghLabel := p.findMatchingGhLabel(allLabels, label)
		if ghLabel == nil {
			if err := p.createLabel(url, label); err != nil {
				return err
			}
			p.log.InfoF("Label '%s' created", label.Name)
			p.delay()
			continue
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
		p.delay()
	}

	return nil
}

// DeleteLabelsForRepo implements core.GitHostingFacade.
func (p *Facade) DeleteLabelsForRepo(url *core.GitURL, labels []core.Label) error {
	p.log.SetName(url.GetRepositoryName())
	converted := LabelConverter{}.ConvertFromEntity(labels)
	for _, label := range converted {
		deleted, err := p.deleteLabel(url, label)
		if err != nil {
			return err
		}
		if deleted {
			p.log.InfoF("Label '%s' deleted", label.Name)
		} else {
			p.log.InfoF("Label '%s' not deleted (not existing)", label.Name)
		}
		p.delay()
	}
	return nil
}

func (p *Facade) createLabel(url *core.GitURL, label *cfg.RepositoryLabel) error {
	_, _, err := p.client.Issues.CreateLabel(p.ctx, url.GetNamespace(), url.GetRepositoryName(), &github.Label{
		Name:        &label.Name,
		Color:       &label.Color,
		Description: &label.Description,
	})
	return err
}

func (p *Facade) updateLabel(url *core.GitURL, ghLabel *github.Label, label *cfg.RepositoryLabel) error {
	// TODO: Without a new_name property we cannot rename a label yet.
	ghLabel.Description = &label.Description
	ghLabel.Color = &label.Color
	_, _, err := p.client.Issues.EditLabel(p.ctx, url.GetNamespace(), url.GetRepositoryName(), label.Name, ghLabel)
	return err
}

func (p *Facade) deleteLabel(url *core.GitURL, label *cfg.RepositoryLabel) (bool, error) {
	resp, err := p.client.Issues.DeleteLabel(p.ctx, url.GetNamespace(), url.GetRepositoryName(), label.Name)
	if resp != nil && resp.StatusCode == 404 {
		// Not an error
		return false, nil
	}
	return err == nil, err
}

func (p *Facade) fetchAllLabels(url *core.GitURL) ([]*github.Label, error) {
	nextPage := 1
	var allLabels []*github.Label
	for repeat := true; repeat; repeat = nextPage > 0 {
		labels, resp, err := p.client.Issues.ListLabels(p.ctx, url.GetNamespace(), url.GetRepositoryName(), &github.ListOptions{
			Page:    nextPage,
			PerPage: 100,
		})
		if err != nil {
			return nil, err
		}
		allLabels = append(allLabels, labels...)
		// On the last page, the NextPage is 0 again, we can use that to exit the loop
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

// delay sleeps one second for abuse rate limit best-practice.
//
// https://docs.github.com/en/rest/guides/best-practices-for-integrators#dealing-with-abuse-rate-limits
// "If you're making a large number of POST, PATCH, PUT, or DELETE requests for a single user or client ID, wait at least one second between each request."
func (p *Facade) delay() {
	time.Sleep(1 * time.Second)
}
