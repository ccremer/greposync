package pullrequest

import (
	"reflect"

	"github.com/ccremer/greposync/core"
	"github.com/ccremer/greposync/printer"
)

// PullRequestHandler contains the business logic to interact with labels on supported core.GitHostingProvider.
type PullRequestHandler struct {
	repoStore     core.GitRepositoryStore
	templateStore core.TemplateStore
	valueStore    core.ValueStore
	log           printer.Printer
}

func NewPullRequestHandler(ts core.TemplateStore, vs core.ValueStore) *PullRequestHandler {
	return &PullRequestHandler{
		templateStore: ts,
		valueStore:    vs,
		log:           printer.New(),
	}
}

func (s *PullRequestHandler) fetchPrTemplate(ctx *pipelineContext) error {
	template, err := s.templateStore.FetchPullRequestTemplate()
	if err != nil {
		return err
	}
	ctx.template = template
	return nil
}

func (s *PullRequestHandler) renderTemplate(ctx *pipelineContext) error {
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

func (s *PullRequestHandler) createOrUpdatePr(ctx *pipelineContext) error {
	if ctx.body != "" {
		ctx.pr.SetBody(ctx.body)
	}
	return ctx.repo.EnsurePullRequest(ctx.pr)
}

func (s *PullRequestHandler) fetchExistingPr(ctx *pipelineContext) error {
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
