package ui

import (
	"bytes"
	"io"
	"os"
	"sync"

	"github.com/pterm/pterm"
)

type ColoredConsole struct {
	// BatchProgressbar is a persistent line appended showing the progress of a batch operation.
	// After each call to Printfln, the line is update.
	// Be sure to assign the new value of this field after calling Start() on the progress bar.
	BatchProgressbar *pterm.ProgressbarPrinter

	buffers map[string]*bytes.Buffer
	m       sync.Mutex

	// Quiet will redirect all console lines to an internal buffer if true.
	Quiet bool
}

func NewColoredConsole() *ColoredConsole {
	return &ColoredConsole{
		BatchProgressbar: pterm.DefaultProgressbar.WithTitle("---------  UPDATING REPOSITORIES..."),
		buffers:          map[string]*bytes.Buffer{},
	}
}

func (c *ColoredConsole) PrintProgressbarMessage(scope string, err error) {
	c.m.Lock()
	defer c.m.Unlock()

	if err == nil {
		pterm.Success.WithScope(pterm.Scope{Text: scope, Style: pterm.Success.Scope.Style}).
			Printfln("Update finished for repository")
	} else {
		pterm.Error.WithScope(pterm.Scope{Text: scope, Style: pterm.Error.Scope.Style}).
			Println("Update failed for repository")
	}
	c.BatchProgressbar.Increment()
}

func (c *ColoredConsole) RefreshProgressBar() {
	c.m.Lock()
	defer c.m.Unlock()

	c.BatchProgressbar.Add(0)
}

func (c *ColoredConsole) AddToBuffer(scope string, buffer io.WriterTo) {
	c.m.Lock()
	defer c.m.Unlock()
	buf := c.getOrCreateBuffer(scope)
	_, _ = buffer.WriteTo(buf)
}

func (c *ColoredConsole) getOrCreateBuffer(scope string) *bytes.Buffer {
	buf, exists := c.buffers[scope]
	if !exists {
		buf = &bytes.Buffer{}
		c.buffers[scope] = buf
	}
	return buf
}

// Flush dumps the logging buffers to stdout.
// This is a noop if the buffers is empty.
func (c *ColoredConsole) Flush(scope, header string) {
	c.m.Lock()
	defer c.m.Unlock()
	buf, exists := c.buffers[scope]
	if exists {
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).Println(header)
		_, _ = buf.WriteTo(os.Stdout)
	}
}
