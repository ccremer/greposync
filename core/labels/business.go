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

func (s *LabelService) createOrUpdateLabels(r core.GitRepository, h core.GitHostingFacade) error {
	labels := filterActiveLabels(r.GetLabels())
	if len(labels) > 0 {
		return h.CreateOrUpdateLabelsForRepo(r.GetConfig().URL, labels)
	}
	return nil
}

func filterActiveLabels(labels []core.Label) []core.Label {
	filtered := make([]core.Label, 0)
	for _, label := range labels {
		if !label.IsBoundForDeletion() {
			filtered = append(filtered, label)
		}
	}
	return filtered
}

func (s *LabelService) deleteLabels(r core.GitRepository, h core.GitHostingFacade) error {
	labels := filterDeadLabels(r.GetLabels())
	if len(labels) > 0 {
		return h.DeleteLabelsForRepo(r.GetConfig().URL, labels)
	}
	return nil
}

func filterDeadLabels(labels []core.Label) []core.Label {
	var filtered []core.Label
	for _, label := range labels {
		if label.IsBoundForDeletion() {
			filtered = append(filtered, label)
		}
	}
	return filtered
}

func (s *LabelService) initHostingAPIs() error {
	occurringProviders := map[core.GitHostingProvider]bool{}
	for _, facade := range s.repoFacades {
		occurringProviders[facade.GetConfig().Provider] = true
	}
	for provider := range occurringProviders {
		if hostingFacade, isSupported := s.repoProvider.GetSupportedGitHostingProviders()[provider]; isSupported {
			if err := hostingFacade.Initialize(); err != nil {
				return err
			}
		} else {
			s.log.WarnF("Provider '%s' is not supported, ignoring all repositories from this Git hosting provider.", provider)
		}
	}
	return nil
}
