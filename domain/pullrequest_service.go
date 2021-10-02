package domain

import pipeline "github.com/ccremer/go-command-pipeline"

type PullRequestService struct {
}

func NewPullRequestService() *PullRequestService {
	return &PullRequestService{}
}

type PullRequestServiceContext struct {
	Repository     *GitRepository
	TemplateEngine TemplateEngine
	Body           string
	Title          string
	TargetBranch   string
}

func (prs *PullRequestService) NewPullRequestForRepository(prsCtx PullRequestServiceContext) error {
	values := Values{
		MetadataValueKey: Values{
			RepositoryValueKey: prsCtx.Repository.AsValues(),
		},
	}

	p := pipeline.NewPipeline().WithSteps(
		pipeline.NewStepFromFunc("body", func(ctx pipeline.Context) error {
			body, err := prsCtx.TemplateEngine.ExecuteString(prsCtx.Body, values)
			prsCtx.Body = body.String()
			return err
		}),
		pipeline.NewStepFromFunc("title", func(ctx pipeline.Context) error {
			title, err := prsCtx.TemplateEngine.ExecuteString(prsCtx.Title, values)
			prsCtx.Title = title.String()
			return err
		}),
		pipeline.NewStepFromFunc("newPR", func(ctx pipeline.Context) error {
			newPr, err := NewPullRequest(nil, prsCtx.Title, prsCtx.Body, prsCtx.Repository.CommitBranch, prsCtx.TargetBranch, prsCtx.Repository.Labels)
			prsCtx.Repository.PullRequest = newPr
			return err
		}),
	)
	return p.Run().Err
}
