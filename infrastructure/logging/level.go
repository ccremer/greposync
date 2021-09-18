package logging

type LogLevel int

const (
	// LevelWarn is the warning logging level.
	LevelWarn = 2
	// LevelInfo is the info logging level.
	LevelInfo = 0
	// LevelDebug is the debug logging level.
	LevelDebug = 1

	LevelSuccess = 3
)

var (
	// DefaultLevel is the default logging level for new logger instances.
	DefaultLevel LogLevel = LevelInfo
)
