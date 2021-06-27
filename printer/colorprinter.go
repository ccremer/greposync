package printer

import (
	"fmt"
	"io"
	"os"
)

type (
	colorPrinter struct {
		name     string
		level    LogLevel
		color    Color
		colorMap map[LogLevel]Color
	}
)

const (
	Black       = "\u001b[30m"
	Red         = "\u001b[31m"
	Green       = "\u001b[32m"
	Yellow      = "\u001b[33m"
	Blue        = "\u001b[34m"
	Magenta     = "\u001b[35m"
	Cyan        = "\u001b[36m"
	White       = "\u001b[37m"
	BrightWhite = "\u001b[37;1m"
	Gray        = "\u001b[38;5;102m"

	Reset = "\u001b[0m"
)

var (
	InfoColor  Color = BrightWhite
	DebugColor Color = Gray
	WarnColor  Color = Yellow
	ErrorColor Color = Red
)

func (p *colorPrinter) SetColor(color Color) {
	p.setColorForWriter(color, os.Stdout)
}

func (p *colorPrinter) UseColor(color Color) Printer {
	p.color = color
	return p
}

func (p *colorPrinter) MapColorToLevel(color Color, level LogLevel) Printer {
	p.colorMap[level] = color
	return p
}

func (p colorPrinter) ResetColor() {
	p.setColorForWriter(Reset, os.Stdout)
}

func (p *colorPrinter) setColorForWriter(color Color, writer io.Writer) Printer {
	_, _ = fmt.Fprint(writer, color)
	return p
}

func (p colorPrinter) DebugF(format string, args ...interface{}) {
	if p.level >= LevelDebug {
		p.printWithColorAndPrefix(os.Stdout, p.colorMap[LevelDebug], format, args...)
	}
}

func (p *colorPrinter) InfoF(format string, args ...interface{}) {
	if p.level >= LevelInfo {
		p.printWithColorAndPrefix(os.Stdout, p.colorMap[LevelInfo], format, args...)
	}
}

func (p *colorPrinter) WarnF(format string, args ...interface{}) {
	if p.level >= LevelWarn {
		p.printWithColorAndPrefix(os.Stdout, p.colorMap[LevelWarn], format, args...)
	}
}

func (p *colorPrinter) LogF(format string, args ...interface{}) {
	p.printWithColor(os.Stdout, p.color, format, args...)
}

func (p *colorPrinter) PrintF(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (p *colorPrinter) getPrefix() string {
	prefix := ""
	if p.name != "" {
		prefix = p.name + ": "
	}
	return prefix
}

func (p *colorPrinter) printWithColorAndPrefix(writer io.Writer, color Color, format string, args ...interface{}) {
	_, _ = fmt.Fprintf(writer, "%s%s%s%s\n", color, p.getPrefix(), fmt.Sprintf(format, args...), Reset)
}

func (p *colorPrinter) printWithColor(writer io.Writer, color Color, format string, args ...interface{}) {
	_, _ = fmt.Fprintf(writer, "%s%s%s\n", color, fmt.Sprintf(format, args...), Reset)
}

func (p *colorPrinter) SetLevel(level LogLevel) Printer {
	p.level = level
	return p
}

func (p *colorPrinter) SetName(name string) Printer {
	p.name = name
	return p
}

func (p *colorPrinter) CheckIfError(err error) {
	if err != nil {
		p.printWithColorAndPrefix(os.Stderr, p.colorMap[LevelError], err.Error())
		os.Exit(1)
	}
}
