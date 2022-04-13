package labels

import (
	"context"

	"github.com/ccremer/greposync/domain"
)

type labelPipeline struct {
	appService     *AppService
	labelsToModify domain.LabelSet
	labelsToDelete domain.LabelSet
	repo           *domain.GitRepository
}

func (p *labelPipeline) updateLabelsForRepository(_ context.Context) error {
	err := p.appService.labelStore.EnsureLabelsForRepository(p.repo, p.repo.Labels)
	return err
}

func (p *labelPipeline) fetchLabelsForRepository(_ context.Context) error {
	labels, err := p.appService.labelStore.FetchLabelsForRepository(p.repo)
	if err != nil {
		return err
	}
	return p.repo.SetLabels(labels)
}

func (p *labelPipeline) deleteLabelsForRepository(_ context.Context) error {
	err := p.appService.labelStore.RemoveLabelsFromRepository(p.repo, p.labelsToDelete)
	return err
}
