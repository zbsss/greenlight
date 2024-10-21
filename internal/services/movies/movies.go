package movies

import (
	"context"

	"github.com/zbsss/greenlight/internal/model"
)

type MovieService struct {
	db *model.Queries
}

func NewMovieService(db *model.Queries) *MovieService {
	return &MovieService{db: db}
}

func (s *MovieService) CreateMovie(ctx context.Context, req CreateMovieRequest) (*MovieResponse, error) {
	// Validate the input
	if err := req.OK(); err != nil {
		return nil, err
	}

	params := model.CreateMovieParams{
		Title:   req.Title,
		Year:    req.Year,
		Runtime: req.Runtime,
		Genres:  req.Genres,
	}

	movie, err := s.db.CreateMovie(ctx, params)
	if err != nil {
		return nil, err
	}

	// Transform the database model to the response model
	return transform(&movie), nil
}

func (s *MovieService) ListMovies(ctx context.Context) ([]*MovieResponse, error) {
	movies, err := s.db.ListMovies(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]*MovieResponse, len(movies))
	for i, movie := range movies {
		response[i] = transform(&movie)
	}
	return response, nil
}

func transform(movie *model.Movie) *MovieResponse {
	return &MovieResponse{
		ID:      movie.ID,
		Title:   movie.Title,
		Year:    movie.Year,
		Runtime: Runtime(movie.Runtime),
		Genres:  movie.Genres,
	}
}
