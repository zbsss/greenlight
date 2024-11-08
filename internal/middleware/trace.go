package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/zbsss/greenlight/internal/contextkeys"
)

const TraceIDHeader = "X-Trace-ID"

func TraceMiddleware(next http.Handler) http.Handler {
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

		ctx := context.WithValue(r.Context(), contextkeys.TraceIDKey, traceID)

		// Set the traceID in the response header
		w.Header().Set(TraceIDHeader, traceID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
