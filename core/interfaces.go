package core

//go:generate rm -rf corefakes
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//CoreService is a representation of a core feature or process.
type CoreService interface {
	// RunPipeline executes the main business logic of this core service.
	// It returns an error if the core service deems the process to have failed.
	RunPipeline() error
}
