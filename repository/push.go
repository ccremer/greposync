package repository

import (
	pipeline "github.com/ccremer/go-command-pipeline"
)

func (s *Service) PushToRemote() pipeline.ActionFunc {
	return func() pipeline.Result {
		args := []string{"push"}
		if s.Config.ForcePush {
			args = append(args, "--force")
		}
		out, stderr, err := s.execGitCommand(s.logArgs(args...)...)
		if err != nil {
			return s.toResult(err, stderr)
		}
		s.p.DebugF(out)
		return pipeline.Result{}
	}
}

func (s *Service) SkipPush() pipeline.Predicate {
	return func(step pipeline.Step) bool {
		return s.Config.SkipPush
	}
}
