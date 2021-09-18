package gotemplate

import "github.com/ccremer/greposync/domain"

type GoTemplateInstrumentation struct {
}

func (i *GoTemplateInstrumentation) fetchedAllTemplatesIfNoError(templates []*domain.Template, err error) {
	if err == nil {

	}
}
