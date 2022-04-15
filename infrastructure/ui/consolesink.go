package ui

import (
	"bytes"
	"os"

	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/ccremer/plogr"
	"github.com/go-logr/logr"
	"github.com/pterm/pterm"
)

// ConsoleSink is a specialized logr.LogSink that uses plogr.PtermSink under the hood, but with extra features like suppression.
type ConsoleSink struct {
	ptermSink *plogr.PtermSink
	console   *ColoredConsole
}

// NewConsoleSink returns a new instance.
// There should only be 1 instance used per application.
func NewConsoleSink(console *ColoredConsole) *ConsoleSink {
	return &ConsoleSink{
		console: console,
	}
}

// Init implements logr.LogSink.
// It will configure log levels that are defined in logging.LogLevel.
func (t *ConsoleSink) Init(info logr.RuntimeInfo) {
	sink := plogr.NewPtermSink()
	sink.Init(info)

	sink.LevelPrinters[logging.LevelDebug] = plogr.DefaultLevelPrinters[1]
	sink.LevelPrinters[logging.LevelInfo] = plogr.DefaultLevelPrinters[0]
	sink.LevelPrinters[logging.LevelSuccess] = pterm.Success
	sink.LevelPrinters[logging.LevelWarn] = pterm.Warning
	sink.ErrorPrinter = *sink.ErrorPrinter.WithLineNumberOffset(3)

	t.ptermSink = &sink
}

// Enabled implements logr.LogSink.
func (t *ConsoleSink) Enabled(level int) bool {
	return t.ptermSink.Enabled(level)
}

// Info implements logr.LogSink.
// If the name is empty or if Quiet is false, the message is always printed to os.Stdout.
// If the name is non-empty, the message will be buffered internally to be printed at once later.
func (t *ConsoleSink) Info(level int, msg string, keysAndValues ...interface{}) {
	buf := &bytes.Buffer{}
	t.ptermSink.WithOutput(buf).Info(level, msg, keysAndValues...)
	if t.ptermSink.Name() == "" || !t.console.Quiet {
		_, _ = buf.WriteTo(os.Stdout)
		return
	}
	t.console.AddToBuffer(t.ptermSink.Name(), buf)
}

// Error implements logr.LogSink.
// If the name is empty or if Quiet is false, the message is always printed to os.Stdout.
// If the name is non-empty, the message will be buffered internally to be printed at once later.
func (t *ConsoleSink) Error(err error, msg string, keysAndValues ...interface{}) {
	buf := &bytes.Buffer{}
	t.ptermSink.WithOutput(buf).Error(err, msg, keysAndValues...)
	if t.ptermSink.Name() == "" || !t.console.Quiet {
		_, _ = buf.WriteTo(os.Stdout)
		return
	}
	t.console.AddToBuffer(t.ptermSink.Name(), buf)
}

// WithValues implements logr.LogSink.
func (t *ConsoleSink) WithValues(keysAndValues ...interface{}) logr.LogSink {
	newSink := &ConsoleSink{
		ptermSink: t.ptermSink.WithValues(keysAndValues...).(*plogr.PtermSink),
		console:   t.console,
	}
	return newSink
}

// WithName implements logr.LogSink.
func (t *ConsoleSink) WithName(name string) logr.LogSink {
	pSink := t.ptermSink.WithName(name).(plogr.PtermSink)
	newSink := &ConsoleSink{
		ptermSink: &pSink,
		console:   t.console,
	}
	return newSink
}

// WithLevelEnabled enables or disables the given log level, if existing.
func (t *ConsoleSink) WithLevelEnabled(level logging.LogLevel, enabled bool) *ConsoleSink {
	t.ptermSink.LevelEnabled[int(level)] = enabled
	return t
}
