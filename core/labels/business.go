package labels

import (
	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/printer"
)

type (
	LabelService struct {
		repoProvider core.ManagedRepoProvider
		repoFacades  []core.GitRepositoryFacade
		log          printer.Printer
	}
)

// NewService returns a new core LabelService instance.
func NewService(repoProvider core.ManagedRepoProvider) *LabelService {
	return &LabelService{
		repoProvider: repoProvider,
		log:          printer.New().SetName("labels"),
	}
}

func (s *LabelService) createOrUpdateLabels(r core.GitRepositoryFacade, h core.GitHostingFacade) error {
	labels := filterActiveLabels(r.GetLabels())
	if len(labels) > 0 {
		return h.CreateOrUpdateLabelsForRepo(r.GetConfig().URL, labels)
	}
	return nil
}

func filterActiveLabels(labels []core.GitRepositoryLabel) []core.GitRepositoryLabel {
	filtered := make([]core.GitRepositoryLabel, 0)
	for _, label := range labels {
		if !label.IsBoundForDeletion() {
			filtered = append(filtered, label)
		}
	}
	return filtered
}

func (s *LabelService) deleteLabels(r core.GitRepositoryFacade, h core.GitHostingFacade) error {
	labels := filterDeadLabels(r.GetLabels())
	if len(labels) > 0 {
		return h.DeleteLabelsForRepo(r.GetConfig().URL, labels)
	}
	return nil
}

func filterDeadLabels(labels []core.GitRepositoryLabel) []core.GitRepositoryLabel {
	var filtered []core.GitRepositoryLabel
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
	for provider, _ := range occurringProviders {
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
