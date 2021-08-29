package update

import (
	"strings"

	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/printer"
)

type Differ struct {
	log        printer.Printer
	repository *domain.GitRepository
}

func (s *Differ) PrettyPrint(diff string) {
	lines := strings.Split(diff, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "-") {
			s.log.UseColor(printer.Red).LogF(line)
			continue
		}
		if strings.HasPrefix(line, "+") {
			s.log.UseColor(printer.Green).LogF(line)
			continue
		}
		s.log.UseColor(printer.White).LogF(line)
	}
}
