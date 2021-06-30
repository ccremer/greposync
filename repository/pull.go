package repository

func (s *Service) Pull() {
	if s.Config.SkipReset {
		s.p.WarnF("Skipped: git pull")
		return
	}
	if s.remoteBranchExists(s.Config.CommitBranch) {
		out, _, err := s.execGitCommand(s.logArgs("pull")...)
		s.p.CheckIfError(err)
		s.p.DebugF(out)
	}
}
