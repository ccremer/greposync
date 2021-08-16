package domain

type LabelStore interface {
	AddLabel(repository *GitRepository, label Label) error
	RemoveLabel(repository *GitRepository, label Label) error
}
