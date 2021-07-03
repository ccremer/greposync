package repository

import (
	"strings"

	pipeline "github.com/ccremer/go-command-pipeline"
)

func (s *Service) SkipCheckoutPredicate() pipeline.Predicate {
	return func(step pipeline.Step) bool {
		return s.Config.SkipReset || s.Config.CommitBranch == ""
	}
}

func (s *Service) CheckoutBranch() pipeline.ActionFunc {
	return func() pipeline.Result {
		out, stderr, err := s.execGitCommand(s.logArgs("checkout", s.Config.CommitBranch)...)
		if err != nil {
			return s.toResult(err, stderr)
		}
		s.p.DebugF(out)
		return pipeline.Result{}
	}
}

func (s *Service) GetDefaultBranch() string {
	out, _, err := s.execGitCommand("remote", "show", "origin")
	s.p.CheckIfError(err)
	lines := strings.Split(out, "\n")
	head := "HEAD branch: "
	for _, line := range lines {
		str := strings.TrimSpace(line)
		if strings.Contains(str, head) {
			return strings.TrimPrefix(str, head)
		}
	}
	s.p.WarnF("No default branch detected. Fall back to master")
	return "master"
}

func (s *Service) localBranchExists(branch string) bool {
	out, _, err := s.execGitCommand("branch", "--list")
	s.p.CheckIfError(err)
	branches := strings.Split(out, "\n")
	for _, line := range branches {
		if strings.Contains(strings.TrimSpace(line), branch) {
			return true
		}
	}
	return false
}

func (s *Service) remoteBranchExists(branch string) bool {
	out, _, err := s.execGitCommand("branch", "-r", "--list")
	s.p.CheckIfError(err)
	branches := strings.Split(out, "\n")
	for _, line := range branches {
		if strings.Contains(strings.TrimSpace(line), branch) {
			return true
		}
	}
	return false
}
