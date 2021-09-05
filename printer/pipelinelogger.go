package printer

import pipeline "github.com/ccremer/go-command-pipeline"

type (
	// PipelineLogger is the implementation for the pipeline lib.
	PipelineLogger struct {
		Logger Printer
	}
)

// Accept prints the name to debug level.
func (p PipelineLogger) Accept(step pipeline.Step) {
	p.Logger.DebugF("executing step '%s'", step.Name)
}
