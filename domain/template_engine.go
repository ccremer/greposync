package domain

type TemplateEngine interface {
	Execute(template *Template, values Values) (RenderResult, error)
}
