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

func (s *LabelStore) FetchLabelsForRepository(url *domain.GitURL) (domain.LabelSet, error) {
	for _, remote := range s.providers {
		if remote.HasSupportFor(url) {
			labels, err := remote.FetchLabels(url)
			return labels, err
		}
	}
	return nil, fmt.Errorf("%s: %w", url, ErrProviderNotSupported)
}

func (s *LabelStore) EnsureLabelsForRepository(url *domain.GitURL, labels domain.LabelSet) error {
	for _, remote := range s.providers {
		if remote.HasSupportFor(url) {
			err := remote.EnsureLabels(url, labels)
			return err
		}
	}
	return fmt.Errorf("%s: %w", url, ErrProviderNotSupported)
}

func (s *LabelStore) RemoveLabelsFromRepository(url *domain.GitURL, labels domain.LabelSet) error {
	for _, remote := range s.providers {
		if remote.HasSupportFor(url) {
			err := remote.DeleteLabels(url, labels)
			return err
		}
	}
	return fmt.Errorf("%s: %w", url, ErrProviderNotSupported)
}
