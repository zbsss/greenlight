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

	"github.com/justinas/alice"
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

func New(cfg Config, handler http.Handler, log *slog.Logger) *Server {
	// common middleware for all APIs
	h := alice.New(
		RecoverPanic,
		rlog.RequestTracingMiddleware(log),
		LogResponseCode,
		SecureHeaders,
	).Then(handler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      h,
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
