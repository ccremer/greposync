package githosting

import (
	"errors"
	"fmt"

	"github.com/ccremer/greposync/domain"
)

type LabelStore struct {
	providers ProviderMap
}

var ErrProviderNotSupported = errors.New("no remote provider found")

func NewLabelStore(providers ProviderMap) *LabelStore {
	return &LabelStore{
		providers: providers,
	}
}

func (s *LabelStore) FetchLabelsForRepository(repository *domain.GitRepository) (domain.LabelSet, error) {
	for _, remote := range s.providers {
		if remote.HasSupportFor(repository.URL) {
			labels, err := remote.FetchLabels(repository)
			return labels, err
		}
	}
	return nil, fmt.Errorf("%s: %w", repository.URL.GetFullName(), ErrProviderNotSupported)
}

func (s *LabelStore) EnsureLabelsForRepository(repository *domain.GitRepository, labels domain.LabelSet) error {
	for _, remote := range s.providers {
		if remote.HasSupportFor(repository.URL) {
			err := remote.EnsureLabels(repository, labels)
			return err
		}
	}
	return fmt.Errorf("%s: %w", repository.URL.GetFullName(), ErrProviderNotSupported)
}

func (s *LabelStore) RemoveLabelsFromRepository(repository *domain.GitRepository, labels domain.LabelSet) error {
	for _, remote := range s.providers {
		if remote.HasSupportFor(repository.URL) {
			err := remote.DeleteLabels(repository, labels)
			return err
		}
	}
	return fmt.Errorf("%s: %w", repository.URL.GetFullName(), ErrProviderNotSupported)
}
