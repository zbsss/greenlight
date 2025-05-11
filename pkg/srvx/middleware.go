package srvx

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func logResponseCode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		Logger(r.Context()).Info(
			"sending response",
			"duration", time.Since(start).String(),
			"statusCode", wrapped.statusCode,
		)
	})
}

func secureHeaders(next http.Handler) http.Handler {
	middleware := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(middleware)
}

func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				ErrServer(w, r, fmt.Errorf("%+v", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

const traceIDHeader = "X-Trace-ID"

func traceRequest(log *slog.Logger) func(http.Handler) http.Handler {
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
