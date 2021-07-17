package repository

import (
	"strings"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/printer"
)

// Diff invokes git to show the changes between HEAD and previous commit.
func (s *Service) Diff() pipeline.ActionFunc {
	return func() pipeline.Result {
		out, stderr, err := s.execGitCommand(s.logArgs("diff", "HEAD~1")...)
		if err != nil {
			if strings.Contains(stderr, "ambiguous argument 'HEAD~1': unknown revision or path not in the working tree.") {
				s.p.InfoF("This is the first commit, no diff available.")
				return pipeline.Result{}
			}
			return s.toResult(err, stderr)
		}
		s.prettyPrint(out)
		return pipeline.Result{}
	}
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
