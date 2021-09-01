package domain

type LabelStore interface {
	FetchLabelsForRepository(url *GitURL) (LabelSet, error)
	EnsureLabelsForRepository(url *GitURL, labels LabelSet) error
	RemoveLabelsFromRepository(url *GitURL, labels LabelSet) error
}
