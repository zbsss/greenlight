package movies

import (
	"context"
	"database/sql"
	"errors"

	"github.com/zbsss/greenlight/internal/movies/model"
)

type MovieService struct {
	db model.Querier
}

func NewMovieService(db model.Querier) *MovieService {
	return &MovieService{db: db}
}

func (s *MovieService) CreateMovie(ctx context.Context, input MovieInput) (*Movie, error) {
	// Validate the input
	if err := input.OK(); err != nil {
		return nil, err
	}

	params := model.CreateMovieParams{
		Title:      input.Title,
		Year:       input.Year,
		RuntimeMin: input.RuntimeMin,
		Genres:     input.Genres,
	}

	movie, err := s.db.CreateMovie(ctx, params)
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

func (s *MovieService) UpdateMovie(ctx context.Context, id int64, input MovieInput) (*Movie, error) {
	if err := input.OK(); err != nil {
		return nil, err
	}

	params := model.UpdateMovieParams{
		ID:         id,
		Title:      input.Title,
		Year:       input.Year,
		RuntimeMin: input.RuntimeMin,
		Genres:     input.Genres,
	}

	movie, err := s.db.UpdateMovie(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMovieNotFound
		}

		return nil, err
	}

	return transform(&movie), nil
}
