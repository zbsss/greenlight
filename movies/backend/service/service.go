package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/zbsss/greenlight/movies/backend/model"
)

var (
	ErrMovieNotFound = errors.New("movie not found")
)

type MovieService struct {
	db model.Querier
}

func New(db model.Querier) *MovieService {
	return &MovieService{db: db}
}

func (s *MovieService) CreateMovie(ctx context.Context, input MovieInput) (*Movie, error) {
	// Validate the input
	if err := input.OK(); err != nil {
		return nil, err
	}

	movie, err := s.db.CreateMovie(ctx, model.CreateMovieParams{
		Title:      input.Title,
		Year:       input.Year,
		RuntimeMin: input.RuntimeMin,
		Genres:     input.Genres,
	})
	if err != nil {
		return nil, err
	}

	// Transform the database model to the response model
	return transform(&movie), nil
}

func (s *MovieService) ListMovies(ctx context.Context) ([]*Movie, error) {
	movies, err := s.db.ListMovies(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]*Movie, len(movies))
	for i, movie := range movies {
		response[i] = transform(&movie)
	}
	return response, nil
}

func (s *MovieService) GetMovie(ctx context.Context, id int64) (*Movie, error) {
	movie, err := s.db.GetMovie(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMovieNotFound
		}

		return nil, err
	}

	return transform(&movie), nil
}

func (s *MovieService) UpdateMovie(ctx context.Context, id int64, updates PartialMovieUpdate) (*Movie, error) {
	movie, err := s.db.GetMovie(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMovieNotFound
		}

		return nil, err
	}

	fullUpdate := mergeMovieUpdates(&movie, &updates)
	if err := fullUpdate.OK(); err != nil {
		return nil, err
	}

	updated, err := s.db.UpdateMovie(ctx, model.UpdateMovieParams{
		ID:         id,
		Title:      fullUpdate.Title,
		Year:       fullUpdate.Year,
		RuntimeMin: fullUpdate.RuntimeMin,
		Genres:     fullUpdate.Genres,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMovieNotFound
		}

		return nil, err
	}

	return transform(&updated), nil
}
