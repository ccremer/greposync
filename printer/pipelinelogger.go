package printer

type (
	// PipelineLogger is the implementation for the pipeline lib.
	PipelineLogger struct {
		Logger Printer
	}
)

// Log prints the message and name to debug level.
func (p PipelineLogger) Log(message, name string) {
	p.Logger.DebugF("%s '%s'", message, name)
}
