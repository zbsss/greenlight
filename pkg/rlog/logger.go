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

var defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

// FromContext returns the Logger instance from the given context.
func FromContext(ctx context.Context) *slog.Logger {
	log, ok := ctx.Value(requestLoggerKey).(*slog.Logger)
	if !ok {
		log = defaultLogger
		log.Warn(
			"Request logger not found in context, using default logger." +
				" Probably RequestTracingMiddleware middleware was not used",
		)
	}
	return log
}

const traceIDHeader = "X-Trace-ID"

func RequestTracingMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var traceID string

			// Check if the incoming request has a X-Trace-ID header
			incomingTraceID := r.Header.Get(traceIDHeader)
			if incomingTraceID != "" {
				traceID = incomingTraceID
			} else {
				// Generate a new trace ID if not present in the request
				traceID = uuid.New().String()
			}

			// Set the traceID in the response header
			w.Header().Set(traceIDHeader, traceID)
			ctx := context.WithValue(r.Context(), traceIDKey, traceID)

			// Create a fresh logger with request details
			requestLog := log.With(
				"traceID", traceID,
				"ip", r.RemoteAddr,
				"proto", r.Proto,
				"method", r.Method,
				"uri", r.URL.RequestURI(),
			)
			ctx = context.WithValue(ctx, requestLoggerKey, requestLog)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
