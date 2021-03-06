package ui

import (
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/go-logr/logr"
)

// ConsoleLoggerFactory creates logr.Logger optimized for console CLI.
type ConsoleLoggerFactory struct {
	rootLogger logr.Logger
	sink       *ConsoleSink
}

// NewConsoleLoggerFactory wraps the given log sink.
func NewConsoleLoggerFactory(sink *ConsoleSink) *ConsoleLoggerFactory {
	return &ConsoleLoggerFactory{
		rootLogger: logr.New(sink),
		sink:       sink,
	}
}

// NewGenericLogger implements logging.LoggerFactory.
func (f *ConsoleLoggerFactory) NewGenericLogger(name string) logr.Logger {
	return f.rootLogger.WithName(name)
}

// NewRepositoryLogger implements logging.LoggerFactory.
func (f *ConsoleLoggerFactory) NewRepositoryLogger(repository *domain.GitRepository) logr.Logger {
	return f.NewGenericLogger(repository.URL.GetFullName())
}

// NewPipelineLogger implements logging.LoggerFactory.
func (f *ConsoleLoggerFactory) NewPipelineLogger(name string) *logging.PipelineLogger {
	return &logging.PipelineLogger{
		Logger: f.NewGenericLogger(name),
	}
}

// SetLogLevel implements logging.LoggerFactory.
func (f *ConsoleLoggerFactory) SetLogLevel(level int) {
	for i := 0; i <= level; i++ {
		f.sink.ptermSink.SetLevelEnabled(i, true)
	}
}
