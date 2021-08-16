package domain

type TemplateRenderer struct {
	engine TemplateEngine
}

func NewTemplateRenderer(engine TemplateEngine) *TemplateRenderer {
	return &TemplateRenderer{engine: engine}
}

func (r *TemplateRenderer) Render(values Values, template *Template) (string, error) {
	return r.engine.Execute(template, values)
}
