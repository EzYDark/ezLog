package ezlog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ezydark/ezlog/log"
	"github.com/fatih/color"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormLogger is a custom logger for Gorm that uses zerolog.
// It should be created using the GormLoggerBuilder.
type GormLogger struct {
	logLevel              logger.LogLevel
	slowThreshold         time.Duration
	sourceField           string
	skipErrRecordNotFound bool
	tag                   string
}

// GormLoggerBuilder is a builder for the GormLogger.
type GormLoggerBuilder struct {
	logger GormLogger
}

// NewGormLogger creates a new GormLoggerBuilder with default values.
func NewGormLogger() *GormLoggerBuilder {
	return &GormLoggerBuilder{
		logger: GormLogger{
			logLevel:              logger.Info, // Default log level
			slowThreshold:         200 * time.Millisecond,
			skipErrRecordNotFound: true,
		},
	}
}

// WithTag adds a custom colored tag to the logger's output.
func (b *GormLoggerBuilder) WithTag(tag string) *GormLoggerBuilder {
	b.logger.tag = tag
	return b
}

// WithLogLevel sets the log level for the logger.
// Valid levels are: Silent, Error, Warn, Info.
func (b *GormLoggerBuilder) WithLogLevel(level logger.LogLevel) *GormLoggerBuilder {
	b.logger.logLevel = level
	return b
}

// WithSlowThreshold sets the slow query threshold.
func (b *GormLoggerBuilder) WithSlowThreshold(threshold time.Duration) *GormLoggerBuilder {
	b.logger.slowThreshold = threshold
	return b
}

// WithSourceField sets the source field for logging.
func (b *GormLoggerBuilder) WithSourceField(field string) *GormLoggerBuilder {
	b.logger.sourceField = field
	return b
}

// WithSkipErrRecordNotFound sets whether to skip gorm.ErrRecordNotFound errors.
func (b *GormLoggerBuilder) WithSkipErrRecordNotFound(skip bool) *GormLoggerBuilder {
	b.logger.skipErrRecordNotFound = skip
	return b
}

// Build creates and returns a configured GormLogger.
func (b *GormLoggerBuilder) Build() *GormLogger {
	return &b.logger
}

// LogMode sets the log mode for the logger.
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// Info logs an info message.
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info {
		log.Info().Msgf(l.formatMsg(msg), data...)
	}
}

// Warn logs a warning message.
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn {
		log.Warn().Msgf(l.formatMsg(msg), data...)
	}
}

// Error logs an error message.
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error {
		log.Error().Msgf(l.formatMsg(msg), data...)
	}
}

// Trace logs a trace message (SQL query).
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	sqlLog := fmt.Sprintf("elapsed=%s rows=%s sql=%s",
		color.New(color.FgYellow).Sprint(elapsed),
		color.New(color.FgCyan).Sprint(rows),
		color.New(color.FgGreen).Sprintf("%q", sql),
	)

	switch {
	case err != nil && (!l.skipErrRecordNotFound || !errors.Is(err, gorm.ErrRecordNotFound)) && l.logLevel >= logger.Error:
		log.Error().Err(err).Msg(l.formatMsg("gorm error " + sqlLog))
	case l.slowThreshold > 0 && elapsed > l.slowThreshold && l.logLevel >= logger.Warn:
		log.Warn().Msg(l.formatMsg("gorm slow query " + sqlLog))
	case l.logLevel >= logger.Info:
		log.Debug().Msg(l.formatMsg("gorm query " + sqlLog))
	}
}

// formatMsg adds the tag to the message if it exists.
func (l *GormLogger) formatMsg(msg string) string {
	if l.tag != "" {
		return fmt.Sprintf("%s %s", color.New(color.FgMagenta).Sprintf("[%s]", l.tag), msg)
	}
	return msg
}
