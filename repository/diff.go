package repository

import (
	"strings"

	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/ccremer/greposync/printer"
)

func (s *Service) ShowDiff() pipeline.ActionFunc {
	return func() pipeline.Result {
		out, stderr, err := s.execGitCommand(s.logArgs("diff", "HEAD~1")...)
		if err != nil {
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
