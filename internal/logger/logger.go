// Package logger provides structured logging.
package logger

import (
	"log/slog"
	"os"
)

// Logger wraps slog for convenient logging.
// All logging methods (Info, Error, Debug, Warn) are available via embedded *slog.Logger.
type Logger struct {
	*slog.Logger
}

// NewLogger creates a new logger instance.
func NewLogger(debug bool) *Logger {
	var level slog.Level
	if debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	baseLogger := slog.New(handler)

	return &Logger{baseLogger}
}
