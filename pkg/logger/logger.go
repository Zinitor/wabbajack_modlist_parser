package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

const LevelFatal = slog.Level(12)

type Interface interface {
	Fatal(message string, args ...any)
	Error(message string, args ...any)
	Warn(message string, args ...any)
	Debug(message string, args ...any)
	Info(message string, args ...any)
}

type Logger struct {
	logger *slog.Logger
}

var _ Interface = (*Logger)(nil)

func New(level string) *Logger {
	var l slog.Level

	switch strings.ToLower(level) {
	case "fatal":
		l = LevelFatal
	case "error":
		l = slog.LevelError
	case "warn":
		l = slog.LevelWarn
	case "info":
		l = slog.LevelInfo
	case "debug":
		l = slog.LevelDebug
	default:
		l = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: l,
		// Add source information (file and line)
		AddSource: true,
		// Replace default attributes
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove the "time" prefix from timestamp to match zerolog format
			if a.Key == slog.TimeKey {
				return slog.Attr{Key: "timestamp", Value: a.Value}
			}
			return a
		},
	}

	// Create JSON handler for structured logging
	handler := slog.NewJSONHandler(os.Stdout, opts)

	// Create the logger
	logger := slog.New(handler)

	// Set as default global logger
	slog.SetDefault(logger)

	return &Logger{
		logger: logger,
	}
}

// Fatal logs a fatal message.
func (l *Logger) Fatal(msg string, args ...any) {
	l.logger.Log(context.Background(), LevelFatal, msg, args...)
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

// Info logs an info message.
func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

// Error logs an error message.
func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

// With returns a new logger with the given attributes.
func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		logger: l.logger.With(args...),
	}
}

// WithError returns a new logger with an "error" attribute.
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		logger: l.logger.With("error", err.Error()),
	}
}
