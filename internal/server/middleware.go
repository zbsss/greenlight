package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zbsss/greenlight/internal/errs"
	"github.com/zbsss/greenlight/internal/rlog"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func LogResponseCode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		rlog.FromContext(r.Context()).Info(
			"sending response",
			"duration", time.Since(start).String(),
			"statusCode", wrapped.statusCode,
		)
	})
}

func SecureHeaders(next http.Handler) http.Handler {
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

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				errs.ServerError(w, r, fmt.Errorf("%+v", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
