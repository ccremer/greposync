package labels

import (
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/printer"
)

type (
	// LabelService contains the business logic to interact with labels on supported core.GitHostingProvider.
	LabelService struct {
		repoProvider core.GitRepositoryStore
		repoFacades  []core.GitRepository
		log          printer.Printer
	}
)

// NewService returns a new core LabelService instance.
func NewService(repoProvider core.GitRepositoryStore) *LabelService {
	return &LabelService{
		repoProvider: repoProvider,
		log:          printer.New().SetName("labels"),
	}
}

func (s *LabelService) createOrUpdateLabels(r core.GitRepository) error {
	labels := r.GetLabels()
	labels = filterActiveLabels(labels)
	if len(labels) <= 0 {
		return nil
	}
	for _, label := range labels {
		changed, err := label.Ensure()
		if err != nil {
			return err
		}
		if changed {
			s.log.InfoF("Label '%s' changed", label.GetName())
		} else {
			s.log.InfoF("Label '%s' unchanged", label.GetName())
		}
	}
	return nil
}

func filterActiveLabels(labels []core.Label) []core.Label {
	filtered := make([]core.Label, 0)
	for _, label := range labels {
		if !label.IsInactive() {
			filtered = append(filtered, label)
		}
	}
	return filtered
}

func (s *LabelService) deleteLabels(r core.GitRepository) error {
	labels := r.GetLabels()
	labels = filterDeadLabels(labels)
	if len(labels) <= 0 {
		return nil
	}
	for _, label := range labels {
		deleted, err := label.Delete()
		if err != nil {
			return err
		}
		if deleted {
			s.log.InfoF("Label '%s' deleted", label.GetName())
		} else {
			s.log.InfoF("Label '%s' not deleted (not existing)", label.GetName())
		}
	}
	return nil
}

func filterDeadLabels(labels []core.Label) []core.Label {
	var filtered []core.Label
	for _, label := range labels {
		if label.IsInactive() {
			filtered = append(filtered, label)
		}
	}
	return filtered
}
