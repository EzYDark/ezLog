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

// LogBuilder is a builder for the global logger.
type LogBuilder struct {
	tviewCompat bool
	writer      io.Writer
}

// New creates a new LogBuilder.
func New() *LogBuilder {
	builder := LogBuilder{
		tviewCompat: false,
		writer:      os.Stdout,
	}
	return &builder
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

// Build builds and configures the zerolog global logger.
func (b *LogBuilder) Build() *zerolog.Logger {
	// Configure global settings
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = "15:04:05.000"

	// Configure custom console writer
	consoleOutput := zerolog.ConsoleWriter{
		Out:        b.writer,
		TimeFormat: "15:04:05.000",
		NoColor:    false,
	}

	consoleOutput.FormatLevel = func(i any) string {
		levelStr := strings.ToUpper(fmt.Sprintf("%s", i))

		if b.tviewCompat {
			// Escape tview's [] text coloring
			switch levelStr {
			case "DEBUG":
				return tview.Escape(color.New(color.FgBlue).Sprintf("[%s]", levelStr))
			case "INFO":
				return tview.Escape(color.New(color.FgGreen).Sprintf("[%s]", levelStr))
			case "WARN":
				return tview.Escape(color.New(color.FgYellow).Sprintf("[%s]", levelStr))
			case "ERROR":
				return tview.Escape(color.New(color.FgRed).Sprintf("[%s]", levelStr))
			case "FATAL":
				return tview.Escape(color.New(color.FgRed, color.Bold).Sprintf("[%s]", levelStr))
			default:
				return tview.Escape(color.New(color.FgWhite).Sprintf("[%s]", levelStr))
			}
		} else {
			switch levelStr {
			case "DEBUG":
				return color.New(color.FgBlue).Sprintf("[%s]", levelStr)
			case "INFO":
				return color.New(color.FgGreen).Sprintf("[%s]", levelStr)
			case "WARN":
				return color.New(color.FgYellow).Sprintf("[%s]", levelStr)
			case "ERROR":
				return color.New(color.FgRed).Sprintf("[%s]", levelStr)
			case "FATAL":
				return color.New(color.FgRed, color.Bold).Sprintf("[%s]", levelStr)
			default:
				return color.New(color.FgWhite).Sprintf("[%s]", levelStr)
			}
		}
	}

	consoleOutput.FormatMessage = func(i any) string {
		return fmt.Sprintf("%s", i)
	}

	consoleOutput.FormatFieldName = func(i any) string {
		return fmt.Sprintf("%s=", i)
	}

	consoleOutput.FormatFieldValue = func(i any) string {
		return fmt.Sprintf("%s", i)
	}

	// Create a new logger instance
	newLogger := zerolog.New(consoleOutput).With().Timestamp().Logger()

	// Set the global logger
	log.Logger = newLogger

	return &newLogger
}
