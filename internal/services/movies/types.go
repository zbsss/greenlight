package movies

import (
	"time"

	"github.com/zbsss/greenlight/internal/model"
	"github.com/zbsss/greenlight/internal/validator"
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
}

type CreateMovieRequest struct {
	Title   string   `json:"title"`
	Year    int32    `json:"year"`
	Runtime int32    `json:"runtime"`
	Genres  []string `json:"genres"`
}

func (m *CreateMovieRequest) OK() error {
	v := validator.New()

	v.Check(m.Title != "", "title", "must be provided")
	v.Check(len(m.Title) <= titleMaxLength, "title", "must not be more than 500 bytes long")

	v.Check(m.Year != 0, "year", "must be provided")
	v.Check(m.Year >= yearMin, "year", "must be greater than 1888")
	v.Check(int(m.Year) <= time.Now().Year(), "year", "must not be in the future")

	v.Check(m.Runtime != 0, "runtime", "must be provided")
	v.Check(m.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(m.Genres != nil, "genres", "must be provided")
	v.Check(len(m.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(m.Genres) <= genresMaxCount, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(m.Genres), "genres", "must not contain duplicate values")

	return v.OK()
}

func transform(movie *model.Movie) *Movie {
	return &Movie{
		ID:      movie.ID,
		Title:   movie.Title,
		Year:    movie.Year,
		Runtime: Runtime(movie.Runtime),
		Genres:  movie.Genres,
	}
}
