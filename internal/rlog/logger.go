package rlog

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type ctxKey string

const (
	traceIDKey       ctxKey = "traceID"
	requestLoggerKey ctxKey = "requestLogger"
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

func (l *Logger) With(args ...any) *Logger {
	return &Logger{l.Logger.With(args...)}
}

// FromContext returns the Logger instance from the given context.
func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value(requestLoggerKey).(*Logger)
	if !ok {
		panic("RequestTracingMiddleware middleware was not used")
	}
	return log
}

const TraceIDHeader = "X-Trace-ID"

func RequestTracingMiddleware(log *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var traceID string

			// Check if the incoming request has a X-Trace-ID header
			incomingTraceID := r.Header.Get(TraceIDHeader)
			if incomingTraceID != "" {
				traceID = incomingTraceID
			} else {
				// Generate a new trace ID if not present in the request
				traceID = uuid.New().String()
			}

			// Set the traceID in the response header
			w.Header().Set(TraceIDHeader, traceID)
			ctx := context.WithValue(r.Context(), traceIDKey, traceID)

			l := log.With("traceID", traceID)
			ctx = context.WithValue(ctx, requestLoggerKey, l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
