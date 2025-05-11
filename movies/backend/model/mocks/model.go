package mocks

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zbsss/greenlight/movies/backend/model"
)

var TestMovie1 = model.Movie{
	ID: 1,
	CreatedAt: pgtype.Timestamptz{
		Time: time.Now(),
	},
	Title:      "Django",
	Year:       2017,
	RuntimeMin: 120,
	Genres: []string{
		"action",
	},
	Version: 1,
}

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

func (mq *MockQueries) Reset(existing ...model.Movie) {
	mq.failOnNext = nil
	mq.movies = map[int64]model.Movie{}

	for _, movie := range existing {
		mq.movies[movie.ID] = movie
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
		Title:      arg.Title,
		Year:       arg.Year,
		RuntimeMin: arg.RuntimeMin,
		Genres:     arg.Genres,
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

func (mq *MockQueries) UpdateMovie(_ context.Context, arg model.UpdateMovieParams) (model.Movie, error) {
	if err := mq.checkForFailure(); err != nil {
		return model.Movie{}, err
	}

	oldMovie, ok := mq.movies[arg.ID]

	if !ok {
		return model.Movie{}, sql.ErrNoRows
	}

	newMovie := model.Movie{
		ID:         arg.ID,
		Version:    oldMovie.Version + 1,
		CreatedAt:  oldMovie.CreatedAt,
		Title:      arg.Title,
		Year:       arg.Year,
		RuntimeMin: arg.RuntimeMin,
		Genres:     arg.Genres,
	}

	mq.movies[newMovie.ID] = newMovie
	return newMovie, nil
}
