package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/zbsss/greenlight/internal/movies/api"
	"github.com/zbsss/greenlight/internal/movies/model"
	movies "github.com/zbsss/greenlight/internal/movies/service"
	"github.com/zbsss/greenlight/internal/server"
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
	bindHealthAPI(app, router)
	api.BindMoviesAPI(app.movies, router)

	srv := server.New(server.Config{Port: cfg.port}, router, logger)

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	return srv.ListenAndShutdownGracefully(ctx)
}

func main() {
	if err := mainNoExit(); err != nil {
		log.Fatalf("%+v", err)
	}
}
