package api

import (
	"fmt"

	"github.com/zbsss/greenlight/movies/backend/service"
)

// toAPIMovie converts a service Movie to an API Movie
func toAPIMovie(serviceMovie *service.Movie) Movie {
	return Movie{
		Id:      serviceMovie.ID,
		Title:   serviceMovie.Title,
		Year:    serviceMovie.Year,
		Runtime: fmt.Sprintf("%d min", serviceMovie.RuntimeMin),
		Genres:  serviceMovie.Genres,
		Version: serviceMovie.Version,
	}
}

func (apiRequest CreateMovieRequest) toService() service.MovieInput {
	return service.MovieInput{
		Title:      apiRequest.Title,
		Year:       apiRequest.Year,
		RuntimeMin: apiRequest.RuntimeMin,
		Genres:     apiRequest.Genres,
	}
}

func (apiRequest UpdateMovieRequest) toService() service.PartialMovieUpdate {
	return service.PartialMovieUpdate{
		Title:      apiRequest.Title,
		Year:       apiRequest.Year,
		RuntimeMin: apiRequest.RuntimeMin,
		Genres:     *apiRequest.Genres,
	}
}
