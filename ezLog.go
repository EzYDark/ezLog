package ezlog

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// globalLogger tracks the global logger instance.
var globalLogger *zerolog.Logger

// LogBuilder is a builder for zerolog loggers.
type LogBuilder struct {
	tviewCompat bool
	writer      io.Writer
	tag         string
	isGlobal    bool
}

// New creates a new LogBuilder, configured by default to create a local logger instance.
func New() *LogBuilder {
	return &LogBuilder{
		tviewCompat: false,
		writer:      os.Stdout,
		isGlobal:    false, // Default behavior is to create a local logger
	}
}

// WithTag adds a custom colored tag to the logger's output.
func (b *LogBuilder) WithTag(tag string) *LogBuilder {
	b.tag = tag
	return b
}

// AsGlobal configures the builder to create a logger that also replaces
// the global instance upon building.
func (b *LogBuilder) AsGlobal() *LogBuilder {
	b.isGlobal = true
	return b
}

// WithTviewCompat sets the tviewCompat field to true.
func (b *LogBuilder) WithTviewCompat() *LogBuilder {
	b.tviewCompat = true
	return b
}

// SetTviewCompat sets the tviewCompat field to true.
func (b *LogBuilder) SetTviewCompat(tviewCompat bool) *LogBuilder {
	b.tviewCompat = tviewCompat
	return b
}

// SetWriter sets the writer field to the given writer.
func (b *LogBuilder) SetWriter(writer io.Writer) *LogBuilder {
	b.writer = writer
	return b
}

// WithWriter sets the writer field to the given writer.
func (b *LogBuilder) WithWriter(writer io.Writer) *LogBuilder {
	b.writer = writer
	return b
}

// Build creates a zerolog.Logger based on the builder's configuration.
func (b *LogBuilder) Build() *zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = "15:04:05.000"

	consoleOutput := zerolog.ConsoleWriter{
		Out:        b.writer,
		TimeFormat: "15:04:05.000",
		NoColor:    false,
	}

	consoleOutput.FormatLevel = func(i any) string {
		levelStr := strings.ToUpper(fmt.Sprintf("%s", i))
		var coloredLevel string

		switch levelStr {
		case "DEBUG":
			coloredLevel = color.New(color.FgBlue).Sprintf("[%s]", levelStr)
		case "INFO":
			coloredLevel = color.New(color.FgGreen).Sprintf("[%s]", levelStr)
		case "WARN":
			coloredLevel = color.New(color.FgYellow).Sprintf("[%s]", levelStr)
		case "ERROR":
			coloredLevel = color.New(color.FgRed).Sprintf("[%s]", levelStr)
		case "FATAL":
			coloredLevel = color.New(color.FgRed, color.Bold).Sprintf("[%s]", levelStr)
		default:
			coloredLevel = color.New(color.FgWhite).Sprintf("[%s]", levelStr)
		}

		if b.tviewCompat {
			return tview.Escape(coloredLevel)
		}
		return coloredLevel
	}

	if b.tag != "" {
		tagStr := color.New(color.FgMagenta).Sprintf("[%s]", b.tag)
		consoleOutput.FormatMessage = func(i any) string {
			return fmt.Sprintf("%s %s", tagStr, i)
		}
	} else {
		consoleOutput.FormatMessage = func(i any) string {
			return fmt.Sprintf("%s", i)
		}
	}

	consoleOutput.FormatFieldName = func(i any) string {
		return color.New(color.FgCyan).Sprintf("%s=", i)
	}

	consoleOutput.FormatFieldValue = func(i any) string {
		if i == nil {
			return color.New(color.FgRed).Sprint("nil")
		}
		switch v := i.(type) {
		case string:
			return color.New(color.FgGreen).Sprintf("%q", v)
		case bool:
			return color.New(color.FgMagenta).Sprint(v)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			return color.New(color.FgYellow).Sprintf("%d", v)
		case float32, float64:
			return color.New(color.FgYellow).Sprintf("%f", v)
		default:
			return fmt.Sprintf("%s", i)
		}
	}

	newLogger := zerolog.New(consoleOutput).With().Timestamp().Logger()

	if b.isGlobal {
		log.Logger = newLogger
		globalLogger = &newLogger
	}

	return &newLogger
}
