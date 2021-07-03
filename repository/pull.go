package repository

import (
	pipeline "github.com/ccremer/go-command-pipeline"
)

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

func (s *Service) SkipReset() pipeline.Predicate {
	return func(step pipeline.Step) bool {
		return s.Config.SkipReset
	}
}
