package printer

type (
	Printer interface {
		// DebugF sets the color to DebugColor, prints the string as given and resets the color, provided the LogLevel is LevelDebug or higher.
		DebugF(format string, args ...interface{})
		// InfoF sets the color to InfoColor, prints the string as given and resets the color, provided the LogLevel is LevelInfo or higher.
		InfoF(format string, args ...interface{})
		// WarnF sets the color to WarnColor, prints the string as given and resets the color, provided the LogLevel is LevelWarn or higher.
		WarnF(format string, args ...interface{})
		// PrintF formats the given format string and prints it to stdout without alterations.
		PrintF(format string, args ...interface{})
		// LogF formats the given format string and prints it to stdout using the color set by UseColor, defaulting to InfoColor if not set.
		LogF(format string, args ...interface{})

		// UseColor will internally save the given color and use it when LogF is invoked.
		UseColor(color Color) Printer
		// MapColorToLevel sets the color for the given log level.
		MapColorToLevel(color Color, level LogLevel) Printer

		// SetColor sets the given color without printing newline char and returns itself.
		// Best used with PrintF as other methods may reset the color before returning.
		SetColor(color Color)
		// ResetColor will reset the color to default.
		// Best used in a defer statement right after setting a color to not forget about resetting.
		ResetColor()
		// SetLevel sets the logging level.
		SetLevel(level LogLevel) Printer
		// CheckIfError will print the error in ErrorColor to stderr if it is non-nil and exit with exit code 1.
		CheckIfError(err error)
	}
	LogLevel int
	Color    string
)

const (
	LevelError = 0
	LevelWarn  = 1
	LevelInfo  = 2
	LevelDebug = 3
)

var (
	DefaultPrinter          = New()
	DefaultLevel   LogLevel = LevelInfo
)

func New() Printer {
	return &colorPrinter{
		level: DefaultLevel,
		color: InfoColor,
		colorMap: map[LogLevel]Color{
			LevelError: ErrorColor,
			LevelWarn:  WarnColor,
			LevelInfo:  InfoColor,
			LevelDebug: DebugColor,
		},
	}
}

func DebugF(format string, args ...interface{}) {
	DefaultPrinter.DebugF(format, args...)
}

func InfoF(format string, args ...interface{}) {
	DefaultPrinter.InfoF(format, args...)
}

func WarnF(format string, args ...interface{}) {
	DefaultPrinter.WarnF(format, args...)
}
func CheckIfError(err error) {
	DefaultPrinter.CheckIfError(err)
}
