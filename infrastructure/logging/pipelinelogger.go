package logging

import (
	pipeline "github.com/ccremer/go-command-pipeline"
	"github.com/go-logr/logr"
)

// PipelineLogger is the implementation for the pipeline lib.
type PipelineLogger struct {
	Logger logr.Logger
}

// Accept prints the scope to debug level.
func (p PipelineLogger) Accept(step pipeline.Step) {
	p.Logger.V(1).Info("Executing step", "step", step.Name)
}
