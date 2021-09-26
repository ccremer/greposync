package ui

import (
	"io"
	"os"
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
	writer io.Writer
	header *pterm.HeaderPrinter
}

// NewConsoleDiffPrinter returns a new instance.
func NewConsoleDiffPrinter() *ConsoleDiffPrinter {
	return &ConsoleDiffPrinter{
		writer: os.Stdout,
		header: pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgMagenta)).WithMargin(15),
	}
}

// PrintDiff implements DiffPrinter.
// The prefix is used to print a header before actually printing the diff.
func (c *ConsoleDiffPrinter) PrintDiff(prefix, diff string) {
	if prefix != "" {
		bytes := c.header.Sprintln(prefix)
		_, _ = c.writer.Write([]byte(bytes))
	}
	lines := strings.Split(diff, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "-") {
			bytes := pterm.FgRed.Sprintln(line)
			_, _ = c.writer.Write([]byte(bytes))
			continue
		}
		if strings.HasPrefix(line, "+") {
			bytes := pterm.FgGreen.Sprintln(line)
			_, _ = c.writer.Write([]byte(bytes))
			continue
		}
		bytes := pterm.Sprintln(line)
		_, _ = c.writer.Write([]byte(bytes))
	}
}
