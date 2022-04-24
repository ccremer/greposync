package ui

import (
	"bytes"
	"io"
	"os"
	"sync"

	"github.com/ccremer/greposync/domain"
	"github.com/mattn/go-isatty"
	"github.com/pterm/pterm"
)

var (
	DefaultProgressbarTitle = "UPDATING REPOSITORIES..."
	titlePrefix             = "---------  "
)

type ColoredConsole struct {
	// batchProgressbar is a persistent line appended showing the progress of a batch operation.
	// After each call to Printfln, the line is updated.
	batchProgressbar *pterm.ProgressbarPrinter

	buffers       map[string]*bytes.Buffer
	m             sync.Mutex
	isInteractive bool

	// Quiet will redirect all console lines to an internal buffer if true.
	Quiet       bool
	commandName string
}

func NewColoredConsole() *ColoredConsole {
	return &ColoredConsole{
		batchProgressbar: pterm.DefaultProgressbar.WithTitle(titlePrefix + DefaultProgressbarTitle),
		buffers:          map[string]*bytes.Buffer{},
		isInteractive:    isatty.IsTerminal(os.Stdout.Fd()),
		commandName:      "Update",
	}
}

func (c *ColoredConsole) StartBatchUpdate(repos []*domain.GitRepository) {
	if !c.isInteractive {
		return
	}
	p, _ := c.batchProgressbar.WithTotal(len(repos)).Start()
	c.batchProgressbar = p
}

func (c *ColoredConsole) PrintProgressbarMessage(scope string, err error) {
	c.m.Lock()
	defer c.m.Unlock()

	if err == nil {
		pterm.Success.WithScope(pterm.Scope{Text: scope, Style: pterm.Success.Scope.Style}).
			Printfln("%s finished for repository", c.commandName)
	} else {
		pterm.Error.WithScope(pterm.Scope{Text: scope, Style: pterm.Error.Scope.Style}).
			Printfln("%s failed for repository", c.commandName)
	}
	if c.isInteractive {
		c.batchProgressbar.Increment()
	}
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

// SetTitle sets the progressbar title.
func (c *ColoredConsole) SetTitle(title string) {
	c.batchProgressbar = c.batchProgressbar.WithTitle(titlePrefix + title)
}

// SetCommandName sets the prefix for success or failure messages
func (c *ColoredConsole) SetCommandName(name string) {
	c.commandName = name
}
