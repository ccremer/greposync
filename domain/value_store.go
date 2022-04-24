package domain

// ValueStore provides methods to query Values from a configuration.
//
// In Domain-Driven Design language, the term `Store` corresponds to `Repository`, but to avoid name clash it was named `Store`.
type ValueStore interface {
	// FetchValuesForTemplate retrieves the Values for the given template.
	FetchValuesForTemplate(template *Template, repository *GitRepository) (Values, error)
	// FetchUnmanagedFlag returns true if the given template should not be rendered.
	// The implementation may return ErrKeyNotFound if the flag is undefined, as the boolean 'false' is ambiguous.
	FetchUnmanagedFlag(template *Template, repository *GitRepository) (bool, error)
	// FetchTargetPath returns an alternative output path for the given template relative to the Git repository.
	// An empty string indicates that there is no alternative path configured.
	FetchTargetPath(template *Template, repository *GitRepository) (Path, error)
	// FetchFilesToDelete returns a slice of Path that should be deleted in the Git repository.
	// The paths are relative to the Git root directory.
	FetchFilesToDelete(repository *GitRepository, templates []*Template) ([]Path, error)
}
