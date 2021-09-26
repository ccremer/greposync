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
func (f *ConsoleLoggerFactory) SetLogLevel(level logging.LogLevel) {
	s := f.sink
	switch level {
	case logging.LevelWarn:
		s.
			WithLevelEnabled(logging.LevelSuccess, false).
			WithLevelEnabled(logging.LevelInfo, false).
			WithLevelEnabled(logging.LevelDebug, false)
	case logging.LevelSuccess:
		s.
			WithLevelEnabled(logging.LevelInfo, false).
			WithLevelEnabled(logging.LevelDebug, false)
	case logging.LevelInfo:
		s.
			WithLevelEnabled(logging.LevelDebug, false)
	}
}
