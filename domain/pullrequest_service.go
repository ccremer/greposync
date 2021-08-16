package domain

type PullRequestService struct {
	prStore PullRequestStore
}

func NewPullRequestService(prStore PullRequestStore) *PullRequestService {
	return &PullRequestService{
		prStore: prStore,
	}
}

// SaveLabelsOnPullRequest persists the labels for the PR in the remote repository.
//
// See also: PullRequestStore.SetLabelsInPullRequest.
func (prs *PullRequestService) SaveLabelsOnPullRequest(repository *GitRepository, pr *PullRequest) error {
	err := prs.prStore.SetLabelsInPullRequest(repository, pr)
	if successful(err) {
		pr.labels = labels
	}
	return err
}

func (prs *PullRequestService) UpdatePullRequestDescription(repository *GitRepository, pr *PullRequest) error {

}
