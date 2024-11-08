package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/zbsss/greenlight/internal/contextkeys"
)

// Logger wraps slog.Logger to provide context-aware logging
type Logger struct {
	*slog.Logger
}

// Create a new logger instance
func NewLogger() *Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// Add function to include caller location in logs
		AddSource: true,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	return &Logger{Logger: logger}
}

// WithContext creates a new logger with trace ID from context
func (l *Logger) WithContext(ctx context.Context) *Logger {
	traceID := ctx.Value(contextkeys.TraceIDKey)
	if traceID == nil {
		return l
	}

	// Create new logger with trace ID added to all logs
	newLogger := l.Logger.With(
		slog.String("traceID", traceID.(string)),
	)

	return &Logger{Logger: newLogger}
}
