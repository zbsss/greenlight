package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/zbsss/greenlight/movies/backend/api"
	"github.com/zbsss/greenlight/movies/backend/service"
	"github.com/zbsss/greenlight/movies/backend/storage"
	"github.com/zbsss/greenlight/pkg/srvx"
)

const defaultPort = 400

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
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

	ms := service.New(storage.New(conn))
	moviesServer := api.NewServer(ms)

	router := http.NewServeMux()
	h := api.HandlerFromMux(moviesServer, router)
	srv := srvx.NewServer(srvx.Config{Port: cfg.port}, h, logger)

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	return srv.ListenAndServe(ctx)
}

func main() {
	if err := mainNoExit(); err != nil {
		log.Fatalf("%+v", err)
	}
}
