package ezlog

import (
	"context"
	"errors"
	"time"

	"github.com/ezydark/ezlog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormLogger is a custom logger for Gorm that uses zerolog.
type GormLogger struct {
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
}

// NewGormLogger creates a new GormLogger.
func NewGormLogger() *GormLogger {
	return &GormLogger{
		SlowThreshold:         200 * time.Millisecond,
		SkipErrRecordNotFound: true,
	}
}

// LogMode sets the log mode for the logger.
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info logs an info message.
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	log.Info().Msgf(msg, data...)
}

// Warn logs a warning message.
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	log.Warn().Msgf(msg, data...)
}

// Error logs an error message.
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	log.Error().Msgf(msg, data...)
}

// Trace logs a trace message (SQL query).
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil && (!l.SkipErrRecordNotFound || !errors.Is(err, gorm.ErrRecordNotFound)) {
		log.Error().Err(err).Str("sql", sql).Dur("elapsed", elapsed).Int64("rows", rows).Msg("gorm error")
		return
	}

	if l.SlowThreshold > 0 && elapsed > l.SlowThreshold {
		log.Warn().Str("sql", sql).Dur("elapsed", elapsed).Int64("rows", rows).Msg("gorm slow query")
		return
	}

	log.Debug().Str("sql", sql).Dur("elapsed", elapsed).Int64("rows", rows).Msg("gorm query")
}
