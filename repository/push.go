package repository

import (
	pipeline "github.com/ccremer/go-command-pipeline"
)

// PushToRemote invokes git to push the commits to origin.
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

// EnabledPush returns true if git pushes are enabled.
func (s *Service) EnabledPush() pipeline.Predicate {
	return func(step pipeline.Step) bool {
		return !s.Config.SkipPush
	}
}
