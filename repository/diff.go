package repository

import (
	"strings"

	"github.com/ccremer/git-repo-sync/printer"
)

func (s *Service) ShowDiff() {
	if s.Config.SkipCommit {
		return
	}
	s.p.DebugF("Getting the latest commit on the current branch")
	ref, err := s.r.Head()
	s.p.CheckIfError(err)

	s.p.DebugF("Retrieving the commit object")
	commit, err := s.r.CommitObject(ref.Hash())
	s.p.CheckIfError(err)

	s.p.DebugF("Retrieving the parent commit")
	parent, err := commit.Parent(0)
	s.p.CheckIfError(err)

	s.p.DebugF("Retrieving the patch between")
	patch, err := parent.Patch(commit)
	s.p.CheckIfError(err)

	s.prettyPrint(patch.String())
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
