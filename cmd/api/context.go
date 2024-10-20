package main

import (
	"log/slog"
	"net/http"
)

type contextKey string

const (
	loggerContextKey  = contextKey("logger")
	traceIDContextKey = contextKey("traceID")
)

func (app *application) loggerFromContext(r *http.Request) *slog.Logger {
	if logger, ok := r.Context().Value(loggerContextKey).(*slog.Logger); ok {
		return logger
	}
	return app.logger
}

func traceIDFromContext(r *http.Request) string {
	return r.Context().Value(traceIDContextKey).(string)
}
