package repository

import (
	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/go-command-pipeline/predicate"
)

// Pull invokes git to pull the latest commits from origin.
func (s *Service) Pull() pipeline.ActionFunc {
	return func() pipeline.Result {
		exists, err := s.remoteBranchExists(s.Config.CommitBranch)
		if err != nil {
			return pipeline.Result{Err: err}
		}
		if exists {
			out, stderr, err := s.execGitCommand(s.logArgs("pull")...)
			if err != nil {
				return s.toResult(err, stderr)
			}
			s.p.DebugF(out)
		}
		return pipeline.Result{}
	}
}

// EnabledReset returns true if git reset is enabled.
func (s *Service) EnabledReset() predicate.Predicate {
	return func(step pipeline.Step) bool {
		return !s.Config.SkipReset
	}
}
