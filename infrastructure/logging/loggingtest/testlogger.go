package loggingtest

import (
	"testing"

	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
)

// TestingLoggerFactory returns a logging.LoggerFactory used for testing.
type TestingLoggerFactory struct {
	t *testing.T
}

// NewTestingLogger returns a new instance with given testing.T.
func NewTestingLogger(t *testing.T) *TestingLoggerFactory {
	return &TestingLoggerFactory{t: t}
}

func (f *TestingLoggerFactory) NewGenericLogger(name string) logr.Logger {
	return testr.New(f.t).WithName(name)
}

func (f *TestingLoggerFactory) NewRepositoryLogger(repository *domain.GitRepository) logr.Logger {
	return f.NewGenericLogger(repository.URL.GetFullName())
}

func (f *TestingLoggerFactory) NewPipelineLogger(name string) *logging.PipelineLogger {
	return &logging.PipelineLogger{Logger: f.NewGenericLogger(name)}
}

func (f *TestingLoggerFactory) SetLogLevel(_ int) {}
