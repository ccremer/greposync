package github

import (
	"time"

	"github.com/ccremer/greposync/core"
	"github.com/google/go-github/v37/github"
)

type LabelImpl struct {
	// Name is the label name.
	Name string `json:"name" koanf:"name"`
	// Description is a short description of the label.
	Description string `json:"description" koanf:"description"`
	// Color is the hexadecimal color code for the label, without the leading #.
	Color string `json:"color" koanf:"color"`
	// Inactive will remove this label.
	Inactive bool `json:"delete" koanf:"delete"`

	remote  *Remote
	repo    *core.GitURL
	ghLabel *github.Label
}

func (l *LabelImpl) IsInactive() bool {
	return l.Inactive
}

func (l *LabelImpl) GetName() string {
	return l.Name
}

func (l *LabelImpl) Delete() (bool, error) {
	return l.remote.deleteLabel(l.repo, l.Name)
}

func (l *LabelImpl) Ensure() (bool, error) {
	if l.ghLabel == nil {
		return true, l.remote.createLabel(l.repo, l)
	}
	if !l.remote.hasLabelChanged(l.ghLabel, l) {
		return false, nil
	}
	return true, l.remote.updateLabel(l.repo, l.ghLabel, l)
}

func (p *Remote) createLabel(url *core.GitURL, label *LabelImpl) error {
	p.m.Lock()
	defer p.delay()
	_, _, err := p.client.Issues.CreateLabel(p.ctx, url.GetNamespace(), url.GetRepositoryName(), &github.Label{
		Name:        &label.Name,
		Color:       &label.Color,
		Description: &label.Description,
	})
	return err
}

func (p *Remote) updateLabel(url *core.GitURL, ghLabel *github.Label, label *LabelImpl) error {
	p.m.Lock()
	defer p.delay()
	// TODO: Without a new_name property we cannot rename a label yet.
	ghLabel.Description = &label.Description
	ghLabel.Color = &label.Color
	_, _, err := p.client.Issues.EditLabel(p.ctx, url.GetNamespace(), url.GetRepositoryName(), label.Name, ghLabel)
	return err
}

func (p *Remote) deleteLabel(url *core.GitURL, name string) (bool, error) {
	p.m.Lock()
	defer p.delay()
	resp, err := p.client.Issues.DeleteLabel(p.ctx, url.GetNamespace(), url.GetRepositoryName(), name)
	if resp != nil && resp.StatusCode == 404 {
		// Not an error
		return false, nil
	}
	return err == nil, err
}

func (p *Remote) fetchAllLabels(url *core.GitURL) ([]*github.Label, error) {
	p.m.Lock()
	defer p.delay()
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

func (p *Remote) findMatchingGhLabel(ghLabels []*github.Label, toFind *LabelImpl) *github.Label {
	for _, label := range ghLabels {
		if label.GetName() == toFind.Name {
			return label
		}
	}
	return nil
}

func (p *Remote) hasLabelChanged(ghLabel *github.Label, repoLabel *LabelImpl) bool {
	return ghLabel.GetDescription() != repoLabel.Description || ghLabel.GetColor() != repoLabel.Color
}

// delay sleeps one second for abuse rate limit best-practice.
//
// https://docs.github.com/en/rest/guides/best-practices-for-integrators#dealing-with-abuse-rate-limits
// "If you're making a large number of POST, PATCH, PUT, or DELETE requests for a single user or client ID, wait at least one second between each request."
func (p *Remote) delay() {
	time.Sleep(1 * time.Second)
	p.m.Unlock()
}
