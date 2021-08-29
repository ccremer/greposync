package domain

type TemplateStore interface {
	FetchTemplates() ([]*Template, error)
}
