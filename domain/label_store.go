package domain

// LabelStore provides methods to interact with labels on a Git hosting service.
//
// In Domain-Driven Design language, the term `Store` corresponds to `Repository`, but to avoid name clash it was named `Store`.
type LabelStore interface {
	// FetchLabelsForRepository retrieves a LabelSet for the given repository.
	FetchLabelsForRepository(repository *GitRepository) (LabelSet, error)
	// EnsureLabelsForRepository creates or updates the given LabelSet in the given repository.
	// Labels that exist remotely, but not in the given LabelSet are ignored.
	// Remote labels have to be updated when Label.GetColor or Label.Description are not matching.
	//
	// Renaming labels are currently not supported.
	EnsureLabelsForRepository(repository *GitRepository, labels LabelSet) error
	// RemoveLabelsFromRepository remotely removes all labels in the given LabelSet.
	// Only the Label.Name is relevant to determine label equality.
	RemoveLabelsFromRepository(repository *GitRepository, labels LabelSet) error
}
