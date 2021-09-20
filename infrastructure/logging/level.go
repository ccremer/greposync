package logging

import (
	"strings"
)

// LogLevel represents the verbosity of application events.
type LogLevel int

const (
	// LevelWarn is the warning logging level.
	LevelWarn = 2
	// LevelInfo is the info logging level.
	LevelInfo = 0
	// LevelDebug is the debug logging level.
	LevelDebug = 1
	// LevelSuccess is the success logging level.
	LevelSuccess = 3
)

// ParseLevelOrDefault returns the parsed LogLevel from a given string.
// If it cannot be parsed, it returns the given default.
func ParseLevelOrDefault(level string, def LogLevel) LogLevel {
	levelMap := map[string]LogLevel{
		"warn":  LevelWarn,
		"info":  LevelInfo,
		"debug": LevelDebug,
	}
	if lvl, found := levelMap[strings.ToLower(level)]; found {
		return lvl
	}
	return def
}
