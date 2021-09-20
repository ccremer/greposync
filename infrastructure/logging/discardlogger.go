package logging

import (
	"github.com/ccremer/greposync/domain"
	"github.com/go-logr/logr"
)

// DiscardLoggerFactory creates logr.Logger that discard everything.
type DiscardLoggerFactory struct{}

// NewDiscardLoggerFactory returns a factory that creates DiscardLogger.
func NewDiscardLoggerFactory() LoggerFactory {
	return &DiscardLoggerFactory{}
}

func (d *DiscardLoggerFactory) NewGenericLogger(_ string) logr.Logger { return logr.Discard() }

func (d *DiscardLoggerFactory) NewRepositoryLogger(_ *domain.GitRepository) logr.Logger {
	return logr.Discard()
}

func (d *DiscardLoggerFactory) NewPipelineLogger(name string) *PipelineLogger {
	return &PipelineLogger{Logger: d.NewGenericLogger(name)}
}
