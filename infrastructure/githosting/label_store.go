package githosting

import (
	"fmt"

	"github.com/ccremer/greposync/domain"
)

type LabelStore struct {
	providers ProviderMap
}

func NewLabelStore(providers ProviderMap) *LabelStore {
	return &LabelStore{
		providers: providers,
	}
}

func (s LabelStore) AddLabel(repository *domain.GitRepository, label domain.Label) error {
	panic("implement me")
}

func (s LabelStore) RemoveLabel(repository *domain.GitRepository, label domain.Label) error {
	panic("implement me")
}

func (s LabelStore) FetchLabelsForRepository(url *domain.GitURL) (domain.LabelSet, error) {
	for _, remote := range s.providers {
		if remote.HasSupportFor(url) {
			labels, err := remote.FindLabels(url)
			return labels, err
		}
	}
	return nil, fmt.Errorf("no remote providers supported: %s", url)
}
