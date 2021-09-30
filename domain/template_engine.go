package domain

// TemplateEngine provides methods to process a Template.
type TemplateEngine interface {
	// Execute renders the given Template with the given Values.
	Execute(template *Template, values Values) (RenderResult, error)
}
