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

// Keys returns a list of keys of the top level.
// Returns an empty string slice if Values is nil or empty.
func (v Values) Keys() []string {
	if v == nil || len(v) == 0 {
		return []string{}
	}
	arr := make([]string, 0)
	for k, _ := range v {
		arr = append(arr, k)
	}
	return arr
}
