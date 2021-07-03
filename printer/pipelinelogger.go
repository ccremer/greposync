package printer

type (
	PipelineLogger struct {
		Logger Printer
	}
)

func (p PipelineLogger) Log(message, name string) {
	p.Logger.DebugF("%s '%s'", message, name)
}
