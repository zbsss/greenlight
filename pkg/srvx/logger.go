package srvx

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

var defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

// Logger returns the Logger instance from the given context.
func Logger(ctx context.Context) *slog.Logger {
	log, ok := ctx.Value(requestLoggerKey).(*slog.Logger)
	if !ok {
		log = defaultLogger
		log.Warn("Request logger not found in context, using default logger.")
	}
	return log
}

// LogErr logs the error to the logger from the given request context.
func LogErr(r *http.Request, err error) {
	Logger(r.Context()).Error(err.Error(), "trace", fmt.Sprintf("%+v", err))
}
