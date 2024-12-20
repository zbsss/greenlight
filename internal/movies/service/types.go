package movies

import (
	"time"

	"github.com/zbsss/greenlight/internal/movies/model"
	"github.com/zbsss/greenlight/pkg/validator"
)

const (
	titleMaxLength = 500
	yearMin        = 1888
	genresMaxCount = 5
)

type Movie struct {
	ID      int64    `json:"id"`
	Title   string   `json:"title"`
	Year    int32    `json:"year"`
	Runtime Runtime  `json:"runtime,omitempty"`
	Genres  []string `json:"genres"`
	Version int32    `json:"version"`
}

type CreateMovieRequest struct {
	Title      string   `json:"title"`
	Year       int32    `json:"year"`
	RuntimeMin int32    `json:"runtimeMin"`
	Genres     []string `json:"genres"`
}

type UpdateMovieRequest struct {
	Title      *string  `json:"title"`
	Year       *int32   `json:"year"`
	RuntimeMin *int32   `json:"runtimeMin"`
	Genres     []string `json:"genres"`
}

// movieInput represents user defined movie data
type movieInput struct {
	Title      string
	Year       int32
	RuntimeMin int32
	Genres     []string
}

func (m movieInput) OK() error {
	v := validator.New()

	v.Check(m.Title != "", "title", "must be provided")
	v.Check(len(m.Title) <= titleMaxLength, "title", "must not be more than 500 bytes long")

	v.Check(m.Year != 0, "year", "must be provided")
	v.Check(m.Year >= yearMin, "year", "must be greater than 1888")
	v.Check(int(m.Year) <= time.Now().Year(), "year", "must not be in the future")

	v.Check(m.RuntimeMin != 0, "runtimeMin", "must be provided")
	v.Check(m.RuntimeMin > 0, "runtimeMin", "must be a positive integer")

	v.Check(m.Genres != nil, "genres", "must be provided")
	v.Check(len(m.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(m.Genres) <= genresMaxCount, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(m.Genres), "genres", "must not contain duplicate values")

	return v.OK()
}

func mergeMovieUpdates(existing *model.Movie, updates *UpdateMovieRequest) movieInput {
	result := movieInput{
		Title:      existing.Title,
		Year:       existing.Year,
		RuntimeMin: existing.RuntimeMin,
		Genres:     existing.Genres,
	}

	if updates.Title != nil {
		result.Title = *updates.Title
	}

	if updates.Year != nil {
		result.Year = *updates.Year
	}

	if updates.RuntimeMin != nil {
		result.RuntimeMin = *updates.RuntimeMin
	}

	if updates.Genres != nil {
		result.Genres = updates.Genres
	}

	return result
}

func transform(movie *model.Movie) *Movie {
	return &Movie{
		ID:      movie.ID,
		Title:   movie.Title,
		Year:    movie.Year,
		Runtime: Runtime(movie.RuntimeMin),
		Genres:  movie.Genres,
		Version: movie.Version,
	}
}
