package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/zbsss/greenlight/internal/errs"
	"github.com/zbsss/greenlight/internal/middleware"
	"github.com/zbsss/greenlight/internal/movies/model"
	movies "github.com/zbsss/greenlight/internal/movies/service"
	"github.com/zbsss/greenlight/internal/rlog"
)

const (
	version         = "0.1.0"
	defaultPort     = 400
	shutdownTimeout = 10 * time.Second
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	log    *slog.Logger
	movies *movies.MovieService
}

func mainNoExit() error {
	var cfg config
	flag.IntVar(&cfg.port, "port", defaultPort, "Port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|prod)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://greenlight:password@localhost/greenlight", "PostgresSQL DSN")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer conn.Close(ctx)

	db := model.New(conn)

	app := &application{
		config: cfg,
		log:    logger,
		movies: movies.NewMovieService(db),
	}

	router := httprouter.New()
	router.MethodNotAllowed = http.HandlerFunc(errs.MethodNotAllowed)
	router.NotFound = http.HandlerFunc(errs.NotFound)

	bindHealthAPI(app, router)
	bindMoviesAPI(app, router)

	// common middleware for all APIs
	handler := alice.New(
		middleware.RecoverPanic,
		rlog.RequestTracingMiddleware(app.log),
		middleware.LogResponseCode,
		middleware.SecureHeaders,
	).Then(router)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				app.log.Info("server shut down gracefully")
			} else {
				app.log.Error("server shut down unexpectedly", "error", err.Error())
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	return srv.Shutdown(ctx)
}

func main() {
	if err := mainNoExit(); err != nil {
		log.Fatalf("%+v", err)
	}
}
