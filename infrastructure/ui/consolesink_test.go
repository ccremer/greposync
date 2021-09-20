package ui

import (
	"testing"
	"time"

	"github.com/go-logr/logr"
)

func TestTerminalSink_Info(t *testing.T) {
	tsink := NewConsoleSink(NewColoredConsole())
	tsink.Init(logr.RuntimeInfo{})

	tsink.Info(0, "message", "key", "value")
	tsink.WithValues("foo", "bar").Info(1, "another", "baz", "fou")
}

func TestTerminalSink_InfoQuiet(t *testing.T) {
	tsink := NewConsoleSink(NewColoredConsole())
	tsink.Init(logr.RuntimeInfo{})

	tsink.console.Quiet = true

	tsink.WithName("repository").Info(0, "message", "key", "value")
	tsink.WithName("repository").WithValues("foo", "bar").Info(1, "another", "baz", "fou")

	tsink.Info(2, "printed message")
	tsink.console.Flush("repository", "header")

}

func TestTerminalSink_InfoParallel(t *testing.T) {
	t.SkipNow()
	tsink := NewConsoleSink(NewColoredConsole())
	tsink.Init(logr.RuntimeInfo{})

	tsink.console.Quiet = true

	go func() {
		tsink.WithName("repository").Info(0, "func 1", "key", "value")
		time.Sleep(1 * time.Second)
		tsink.WithName("repository").WithValues("f", 1).Info(1, "another", "baz", "fou")
	}()

	go func() {
		time.Sleep(500 * time.Millisecond)
		tsink.WithName("repository").Info(0, "func 2", "key", "value")
		time.Sleep(1 * time.Second)
		tsink.WithName("repository").WithValues("f", 2).Info(1, "another", "baz", "fou")
	}()

	tsink.Info(2, "printed message")
	time.Sleep(4 * time.Second)
	tsink.console.Flush("repository", "header")
}
