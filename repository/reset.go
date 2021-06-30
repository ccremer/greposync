package repository

func (s *Service) ResetRepository() {
	if s.Config.SkipReset {
		s.p.WarnF("Skipped: git reset")
		return
	}
	out, _, err := s.execGitCommand(s.logArgs("fetch")...)
	s.p.CheckIfError(err)
	if out != "" {
		s.p.InfoF(out)
	}

	out, _, err = s.execGitCommand(s.logArgs("reset", "--hard")...)
	s.p.CheckIfError(err)
	s.p.DebugF(out)
}
