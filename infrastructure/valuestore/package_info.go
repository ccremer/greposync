package valuestore

var (
	SyncConfigFileName   = ".sync.yml"
	GlobalConfigFileName = "config_defaults.yml"
)

// config is just an alias for easier readability.
type config map[string]interface{}
