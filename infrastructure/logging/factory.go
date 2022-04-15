package logging

import (
	"github.com/ccremer/greposync/domain"
	"github.com/go-logr/logr"
)

type LoggerFactory interface {
	NewGenericLogger(name string) logr.Logger
	NewRepositoryLogger(repository *domain.GitRepository) logr.Logger
	NewPipelineLogger(name string) *PipelineLogger
	SetLogLevel(level int)
}
