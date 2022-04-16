package ui

import (
	"strings"

	"github.com/pterm/pterm"
)

// DiffPrinter is optimized for printing diff from Git output
type DiffPrinter interface {
	// PrintDiff prints the diff.
	// The prefix can be used to identify which scope this diff belongs to.
	PrintDiff(prefix string, diff string)
}

// ConsoleDiffPrinter prints a colored diff to console.
type ConsoleDiffPrinter struct {
	header *pterm.HeaderPrinter
}

// NewConsoleDiffPrinter returns a new instance.
func NewConsoleDiffPrinter() *ConsoleDiffPrinter {
	return &ConsoleDiffPrinter{
		header: pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgMagenta)).WithMargin(15),
	}
}

// PrintDiff implements DiffPrinter.
// The prefix is used to print a header before actually printing the diff.
func (c *ConsoleDiffPrinter) PrintDiff(prefix, diff string) {
	lines := strings.Split(diff, "\n")
	if diff == "" {
		// no diff to print
		return
	}
	if prefix != "" {
		c.header.Println(prefix)
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "-") {
			pterm.FgRed.Println(line)
			continue
		}
		if strings.HasPrefix(line, "+") {
			pterm.FgGreen.Println(line)
			continue
		}
		pterm.Println(line)
	}
}
