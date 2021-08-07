package pullrequest

import (
	"reflect"

	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/printer"
)

// PullRequestService contains the business logic to interact with labels on supported core.GitHostingProvider.
type PullRequestService struct {
	templateStore core.TemplateStore
	valueStore    core.ValueStore
	log           printer.Printer
}


func NewInstance() *PullRequestService {
	return &PullRequestService{
		templateStore: nil,
		valueStore:    nil,
		log:           printer.New(),
	}
}

func (s *PullRequestService) fetchPrTemplate(ctx *pipelineContext) error {
	template, err := s.templateStore.FetchPullRequestTemplate()
	if err != nil {
		return err
	}
	ctx.template = template
	return nil
}

func (s *PullRequestService) renderTemplate(ctx *pipelineContext) error {
	if isNil(ctx.template) {
		return nil
	}
	body, err := ctx.template.Render(core.Values{
		"Metadata":    ctx.repo.GetConfig(),
		"PullRequest": ctx.pr,
	})
	if err != nil {
		return err
	}
	ctx.body = body
	return nil
}

func (s *PullRequestService) createOrUpdatePr(ctx *pipelineContext) error {
	if ctx.body != "" {
		ctx.pr.SetBody(ctx.body)
	}
	return ctx.repo.EnsurePullRequest(ctx.pr)
}

func (s *PullRequestService) fetchExistingPr(ctx *pipelineContext) error {
	pr, err := ctx.repo.FetchPullRequest()
	if err != nil {
		return err
	}
	if isNil(pr) {
		ctx.pr = ctx.repo.NewPullRequest()
	}
	return nil
}

func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}
