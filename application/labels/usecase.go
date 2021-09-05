package labels

import (
	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/domain"
)

type labelUseCase struct {
	appService     *AppService
	labelsToModify domain.LabelSet
	labelsToDelete domain.LabelSet
}

func (uc *labelUseCase) updateLabelsForRepositoryAction(r *domain.GitRepository) pipeline.ActionFunc {
	return func(ctx pipeline.Context) pipeline.Result {
		err := uc.updateLabelsForRepository(r)
		return pipeline.Result{Err: err}
	}
}

func (uc *labelUseCase) updateLabelsForRepository(r *domain.GitRepository) error {
	err := uc.appService.labelStore.EnsureLabelsForRepository(r.URL, r.Labels)
	return err
}

func (uc *labelUseCase) fetchLabelsForRepository(r *domain.GitRepository) pipeline.ActionFunc {
	return func(_ pipeline.Context) pipeline.Result {
		labels, err := uc.appService.labelStore.FetchLabelsForRepository(r.URL)
		if err != nil {
			return pipeline.Result{Err: err}
		}
		err = r.SetLabels(labels)
		return pipeline.Result{Err: err}
	}
}

func (uc *labelUseCase) deleteLabelsForRepository(r *domain.GitRepository) pipeline.ActionFunc {
	return func(ctx pipeline.Context) pipeline.Result {
		err := uc.appService.labelStore.RemoveLabelsFromRepository(r.URL, uc.labelsToDelete)
		return pipeline.Result{Err: err}
	}
}
