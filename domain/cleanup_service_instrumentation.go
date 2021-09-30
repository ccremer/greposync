package domain

// CleanupServiceInstrumentation provides methods for domain observability.
type CleanupServiceInstrumentation interface {
	// FetchedFilesToDelete logs a message indicating that fetching file paths to delete from ValueStore was successful but only if fetchErr is nil.
	// Returns fetchErr unmodified for method chaining.
	FetchedFilesToDelete(fetchErr error, files []Path) error
	// DeletedFile logs a message indicating that deleting file occurred.
	DeletedFile(file Path)
	// WithRepository returns an instance that has the given repository as scope.
	WithRepository(repository *GitRepository) CleanupServiceInstrumentation
}
