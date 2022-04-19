package loggingtest

import (
	"github.com/ccremer/greposync/domain"
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/go-logr/logr"
)

// DiscardLoggerFactory creates logr.Logger that discard everything.
type DiscardLoggerFactory struct{}

// NewDiscardLoggerFactory returns a factory that creates DiscardLogger.
func NewDiscardLoggerFactory() logging.LoggerFactory {
	return &DiscardLoggerFactory{}
}

func (d *DiscardLoggerFactory) SetLogLevel(_ int) {}

func (d *DiscardLoggerFactory) NewGenericLogger(_ string) logr.Logger { return logr.Discard() }

func (d *DiscardLoggerFactory) NewRepositoryLogger(_ *domain.GitRepository) logr.Logger {
	return logr.Discard()
}

func (d *DiscardLoggerFactory) NewPipelineLogger(name string) *logging.PipelineLogger {
	return &logging.PipelineLogger{Logger: d.NewGenericLogger(name)}
}
