package domain

// TemplateStore provides methods to load Template from template root directory.
//
// In Domain-Driven Design language, the term `Store` corresponds to `Repository`, but to avoid name clash it was named `Store`.
type TemplateStore interface {
	// FetchTemplates lists all templates.
	// It aborts on first error.
	FetchTemplates() ([]*Template, error)
}
