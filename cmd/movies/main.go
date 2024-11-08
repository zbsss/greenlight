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
	"github.com/zbsss/greenlight/internal/logger"
	"github.com/zbsss/greenlight/internal/movies/model"
	movies "github.com/zbsss/greenlight/internal/movies/service"
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
	log    *logger.Logger
	movies *movies.MovieService
}

func mainNoExit() error {
	var cfg config
	flag.IntVar(&cfg.port, "port", defaultPort, "Port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|prod)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://greenlight:password@localhost/greenlight", "PostgresSQL DSN")
	flag.Parse()

	mlog := logger.NewLogger()

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		mlog.Error(err.Error())
		os.Exit(1)
	}

	defer conn.Close(ctx)

	db := model.New(conn)

	app := application{
		config: cfg,
		log:    mlog,
		movies: movies.NewMovieService(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(mlog.Handler(), slog.LevelError),
	}

	mlog.Info("starting server", "addr", srv.Addr, "env", cfg.env)

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
