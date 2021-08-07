package core

// Label is attached to a remote Git repository on a supported Git hosting provider.
// The implementation may contain additional provider-specific properties.
//counterfeiter:generate . Label
type Label interface {
	// IsInactive returns true if the label is bound for removal from a remote repository.
	IsInactive() bool
	// GetName returns the label name.
	GetName() string

	// Delete removes the label from the remote repository.
	Delete() (bool, error)
	// Ensure creates the label in the remote repository if it doesn't exist.
	// If the label already exists, it will be updated if the properties are different.
	Ensure() (bool, error)
}
