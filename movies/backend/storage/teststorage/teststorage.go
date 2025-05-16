package teststorage

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/jackc/pgx/v5"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/zbsss/greenlight/movies/backend/storage"
	"github.com/zbsss/greenlight/movies/backend/storage/migrations"
)

type TestStorage struct {
	storage.Querier
	container *postgres.PostgresContainer
}

func New(ctx context.Context) (*TestStorage, error) {
	pg, err := newContainer(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run postgres")
	}

	connectionString, err := pg.ConnectionString(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get connection string")
	}

	if err := migrations.Up(connectionString + "sslmode=disable"); err != nil {
		return nil, errors.Wrap(err, "failed to run migrations")
	}

	conn, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to postgres")
	}

	q := storage.New(conn)

	if err := seedMockData(ctx, q); err != nil {
		return nil, errors.Wrap(err, "failed to seed mock data")
	}

	return &TestStorage{q, pg}, nil
}

func (ts *TestStorage) Close(ctx context.Context) error {
	return ts.container.Terminate(ctx)
}

func newContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	pg, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithUsername("greenlight"),
		postgres.WithPassword("password"),
		postgres.WithDatabase("greenlight"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, err
	}

	return pg, nil
}

func seedMockData(ctx context.Context, q storage.Querier) error {
	movies := []struct {
		title      string
		year       int32
		runtimeMin int32
		genres     []string
	}{
		{
			title:      "The Shawshank Redemption",
			year:       1994,
			runtimeMin: 142,
			genres:     []string{"drama"},
		},
		{
			title:      "The Godfather",
			year:       1972,
			runtimeMin: 175,
			genres:     []string{"crime", "drama"},
		},
		{
			title:      "Pulp Fiction",
			year:       1994,
			runtimeMin: 154,
			genres:     []string{"crime", "drama"},
		},
		{
			title:      "The Dark Knight",
			year:       2008,
			runtimeMin: 152,
			genres:     []string{"action", "crime", "drama"},
		},
		{
			title:      "Fight Club",
			year:       1999,
			runtimeMin: 139,
			genres:     []string{"drama"},
		},
	}

	for _, m := range movies {
		_, err := q.CreateMovie(ctx, storage.CreateMovieParams{
			Title:      m.title,
			Year:       m.year,
			RuntimeMin: m.runtimeMin,
			Genres:     m.genres,
		})
		if err != nil {
			return err
		}

	}
	ms, err := q.ListMovies(ctx)
	if err != nil {
		return err
	}
	if len(ms) != 5 {
		return fmt.Errorf("expected 5 movies, got %d", len(ms))
	}

	return nil
}
