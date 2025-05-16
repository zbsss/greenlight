package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/zbsss/greenlight/movies/backend/api"
	"github.com/zbsss/greenlight/movies/backend/service"
	"github.com/zbsss/greenlight/movies/backend/storage"
	"github.com/zbsss/greenlight/movies/backend/storage/teststorage"

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
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgresSQL DSN")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx := context.Background()

	movieStorage, cleanup, err := setupStorage(ctx, cfg.env, cfg.db.dsn)
	if err != nil {
		return err
	}
	defer cleanup(ctx)

	ms := service.New(movieStorage)
	moviesServer := api.NewServer(ms)

	router := http.NewServeMux()
	h := api.HandlerFromMux(moviesServer, router)
	srv := srvx.NewServer(srvx.Config{Port: cfg.port}, h, logger)

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	return srv.ListenAndServe(ctx)
}

func setupStorage(ctx context.Context, env string, dsn string) (storage.Querier, func(context.Context) error, error) {
	if env == "dev" && dsn == "" {
		ts, err := teststorage.New(ctx)
		if err != nil {
			return nil, nil, err
		}
		return ts, ts.Close, nil
	} else if env == "prod" {
		conn, err := pgx.Connect(ctx, dsn)
		if err != nil {
			return nil, nil, err
		}
		return storage.New(conn), func(ctx context.Context) error {
			return conn.Close(ctx)
		}, nil
	}
	return nil, nil, fmt.Errorf("unsupported environment: %s", env)
}

func main() {
	if err := mainNoExit(); err != nil {
		log.Fatalf("%+v", err)
	}
}
