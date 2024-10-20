package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func secureHeaders(next http.Handler) http.Handler {
	middleware := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(middleware)
}

const traceIDHeader = "X-Trace-ID"

func withTraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var traceID string

		// Check if the incoming request has a Trace-ID header
		incomingTraceID := r.Header.Get(traceIDHeader)
		if incomingTraceID != "" {
			traceID = incomingTraceID
		} else {
			// Generate a new trace ID if not present in the request
			traceID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), traceIDContextKey, traceID)

		// Set the traceID in the response header
		w.Header().Set(traceIDHeader, traceID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			traceID = traceIDFromContext(r)
			ip      = r.RemoteAddr
			proto   = r.Proto
			method  = r.Method
			uri     = r.URL.RequestURI()
		)

		logger := app.logger.With("traceID", traceID, "ip", ip, "proto", proto, "method", method, "uri", uri)
		logger.Info("received request")

		// Add the logger to the request context
		ctx := context.WithValue(r.Context(), loggerContextKey, logger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.serverError(w, r, fmt.Errorf("%+v", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
