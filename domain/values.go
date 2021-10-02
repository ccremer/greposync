package domain

// Values contain a tree of properties to be consumed by a TemplateEngine.
type Values map[string]interface{}

const (
	// MetadataValueKey is the root key for the metadata variables.
	MetadataValueKey = "Metadata"
	// RepositoryValueKey is the key for the GitRepository variable.
	RepositoryValueKey = "Repository"
	// TemplateValueKey is the key for the Template variable.
	TemplateValueKey = "Template"
	// ValuesKey is the key for user-defined variables.
	ValuesKey = "Values"
)
