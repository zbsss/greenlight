package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/zbsss/greenlight/internal/movies/api"
	"github.com/zbsss/greenlight/internal/movies/model"
	movies "github.com/zbsss/greenlight/internal/movies/service"
)

const (
	version     = "0.1.0"
	defaultPort = 400
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
		return err
	}

	defer conn.Close(ctx)

	db := model.New(conn)
	ms := movies.NewMovieService(db)

	// app := &application{
	// 	config: cfg,
	// 	log:    logger,
	// 	movies: movies.NewMovieService(db),
	// }

	router := http.NewServeMux()
	moviesServer := api.NewServer(ms)

	h := api.HandlerFromMux(moviesServer, router)

	srv := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:8080",
	}

	// srv := server.New(server.Config{Port: cfg.port}, router, logger)

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	return srv.ListenAndServe()
	// return srv.ListenAndShutdownGracefully(ctx)
	// return nil
}

func main() {
	if err := mainNoExit(); err != nil {
		log.Fatalf("%+v", err)
	}
}
