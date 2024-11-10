package mocks

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zbsss/greenlight/internal/movies/model"
)

type MockQueries struct {
	movies     map[int64]model.Movie
	failOnNext error
}

var _ model.Querier = &MockQueries{}

func NewMockQueries() *MockQueries {
	mq := &MockQueries{}
	mq.Reset()
	return mq
}

func (mq *MockQueries) Reset() {
	mq.failOnNext = nil
	mq.movies = map[int64]model.Movie{
		1: {
			ID: 1,
			CreatedAt: pgtype.Timestamptz{
				Time: time.Now(),
			},
			Title:   "Django",
			Year:    2017,
			Runtime: 120,
			Genres: []string{
				"action",
			},
			Version: 1,
		},
	}
}

func (mq *MockQueries) FailOnNextCall(err error) {
	mq.failOnNext = err
}

func (mq *MockQueries) checkForFailure() error {
	if mq.failOnNext != nil {
		err := mq.failOnNext
		mq.failOnNext = nil
		return err
	}

	return nil
}

func (mq *MockQueries) CreateMovie(_ context.Context, arg model.CreateMovieParams) (model.Movie, error) {
	if err := mq.checkForFailure(); err != nil {
		return model.Movie{}, err
	}

	movie := model.Movie{
		ID:      int64(len(mq.movies)) + 1,
		Version: 1,
		CreatedAt: pgtype.Timestamptz{
			Time: time.Now(),
		},
		Title:   arg.Title,
		Year:    arg.Year,
		Runtime: arg.Runtime,
		Genres:  arg.Genres,
	}

	mq.movies[movie.ID] = movie

	return movie, nil
}

func (mq *MockQueries) GetMovie(_ context.Context, id int64) (model.Movie, error) {
	if err := mq.checkForFailure(); err != nil {
		return model.Movie{}, err
	}

	movie, ok := mq.movies[id]

	if !ok {
		return model.Movie{}, sql.ErrNoRows
	}

	return movie, nil
}

func (mq *MockQueries) ListMovies(_ context.Context) ([]model.Movie, error) {
	if err := mq.checkForFailure(); err != nil {
		return nil, err
	}

	movies := make([]model.Movie, len(mq.movies))
	i := 0
	for _, movie := range mq.movies {
		movies[i] = movie
		i++
	}

	return movies, nil
}
