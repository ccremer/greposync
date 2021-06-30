package repository

func (s *Service) PushToRemote() {
	if s.Config.SkipPush {
		s.p.WarnF("Skipped: git push")
		return
	}
	args := []string{"push"}
	if s.Config.ForcePush {
		args = append(args, "--force")
	}
	out, _, err := s.execGitCommand(s.logArgs(args...)...)
	s.p.DebugF(out)
	s.p.CheckIfError(err)
}
