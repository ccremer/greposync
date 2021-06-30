package repository

import (
	"strings"

	"github.com/ccremer/git-repo-sync/printer"
)

func (s *Service) ShowDiff() {
	if s.Config.SkipCommit {
		return
	}
	out, _, err := s.execGitCommand(s.logArgs("diff", "HEAD~1")...)
	s.p.CheckIfError(err)
	s.prettyPrint(out)
}

func (s *Service) prettyPrint(diff string) {
	lines := strings.Split(diff, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "-") {
			s.p.UseColor(printer.Red).LogF(line)
			continue
		}
		if strings.HasPrefix(line, "+") {
			s.p.UseColor(printer.Green).LogF(line)
			continue
		}
		s.p.UseColor(printer.White).LogF(line)
	}
}
