package domain

// RenderServiceInstrumentationFactory provides factory methods to create new instances.
type RenderServiceInstrumentationFactory interface {
	// NewRenderServiceInstrumentation creates a new RenderServiceInstrumentation instance using the given GitRepository as context.
	NewRenderServiceInstrumentation(repository *GitRepository) RenderServiceInstrumentation
}

// RenderServiceInstrumentation provides methods for domain observability.
type RenderServiceInstrumentation interface {
	// FetchedTemplatesFromStore logs a message indicating that fetching templates from TemplateStore was successful, but only if fetchErr is nil.
	// Returns fetchErr unmodified for method chaining.
	FetchedTemplatesFromStore(fetchErr error) error
	// FetchedValuesForTemplate logs a message indicating that fetching Values from ValueStore was successful but only if fetchErr is nil.
	// Returns fetchErr unmodified for method chaining.
	FetchedValuesForTemplate(fetchErr error, template *Template) error
	// AttemptingToRenderTemplate logs a message indicating that the actual rendering is about to begin.
	AttemptingToRenderTemplate(template *Template)
	WrittenRenderResultToFile(template *Template, targetPath Path, writeErr error) error
}
