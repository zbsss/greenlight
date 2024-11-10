package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/zbsss/greenlight/pkg/errs"
	"github.com/zbsss/greenlight/pkg/rlog"
)

const (
	shutdownTimeout = 10 * time.Second
)

type Config struct {
	Port int
}

type Server struct {
	*http.Server
	log *slog.Logger
}

func New(cfg Config, router *httprouter.Router, log *slog.Logger) *Server {
	router.MethodNotAllowed = http.HandlerFunc(errs.MethodNotAllowed)
	router.NotFound = http.HandlerFunc(errs.NotFound)
	router.HandlerFunc("GET", "/info/health", healthcheck)

	// common middleware for all APIs
	handler := alice.New(
		RecoverPanic,
		rlog.RequestTracingMiddleware(log),
		LogResponseCode,
		SecureHeaders,
	).Then(router)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(log.Handler(), slog.LevelError),
	}

	return &Server{srv, log}
}

func (s *Server) ListenAndShutdownGracefully(ctx context.Context) error {
	go func() {
		if err := s.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				s.log.Info("server shut down gracefully")
			} else {
				s.log.Error("server shut down unexpectedly", "error", err.Error())
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	return s.Shutdown(ctx)
}
