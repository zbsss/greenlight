package api

import (
	"errors"
	"fmt"
	"net/http"

	movies "github.com/zbsss/greenlight/internal/movies/service"
	"github.com/zbsss/greenlight/pkg/body"
	"github.com/zbsss/greenlight/pkg/errs"
	"github.com/zbsss/greenlight/pkg/rlog"
	"github.com/zbsss/greenlight/pkg/validator"
)

type Server struct {
	ms *movies.MovieService
}

func NewServer(ms *movies.MovieService) Server {
	return Server{ms: ms}
}

func (s Server) GetV1Movies(w http.ResponseWriter, r *http.Request) {
	mvs, err := s.ms.ListMovies(r.Context())
	if err != nil {
		errs.ServerError(w, r, err)
		return
	}

	err = body.WriteJSON(w, http.StatusOK, body.Envelope{"movies": mvs}, nil)
	if err != nil {
		errs.ServerError(w, r, err)
		return
	}
}

func (s Server) PostV1Movies(w http.ResponseWriter, r *http.Request) {
	var apiInput CreateMovieRequest
	err := body.ReadJSON(w, r, &apiInput)
	if err != nil {
		errs.BadRequest(w, r, err)
		return
	}

	serviceInput := movies.CreateMovieRequest{
		Title:      apiInput.Title,
		Year:       apiInput.Year,
		RuntimeMin: apiInput.RuntimeMin,
		Genres:     apiInput.Genres,
	}

	movie, err := s.ms.CreateMovie(r.Context(), serviceInput)
	if err != nil {
		var validationErr validator.ValidationError
		if errors.As(err, &validationErr) {
			errs.BadRequest(w, r, err)
			return
		}

		errs.ServerError(w, r, err)
		return
	}

	rlog.FromContext(r.Context()).Info("created movie", "movie", movie)

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = body.WriteJSON(w, http.StatusCreated, body.Envelope{"movie": movie}, headers)
	if err != nil {
		errs.ServerError(w, r, err)
		return
	}
}

func (s Server) GetV1MoviesId(w http.ResponseWriter, r *http.Request, id int64) {
	movie, err := s.ms.GetMovie(r.Context(), id)
	if err != nil {
		if errors.Is(err, movies.ErrMovieNotFound) {
			errs.NotFound(w, r)
			return
		}

		errs.ServerError(w, r, err)
		return
	}

	err = body.WriteJSON(w, http.StatusOK, body.Envelope{"movie": movie}, nil)
	if err != nil {
		errs.ServerError(w, r, err)
		return
	}
}

func (s Server) PatchV1MoviesId(w http.ResponseWriter, r *http.Request, id int64) {
	var apiInput UpdateMovieRequest
	err := body.ReadJSON(w, r, &apiInput)
	if err != nil {
		errs.BadRequest(w, r, err)
		return
	}

	serviceInput := movies.UpdateMovieRequest{
		Title:      apiInput.Title,
		Year:       apiInput.Year,
		RuntimeMin: apiInput.RuntimeMin,
		Genres:     *apiInput.Genres,
	}

	movie, err := s.ms.UpdateMovie(r.Context(), id, serviceInput)
	if err != nil {
		var validationErr validator.ValidationError
		if errors.As(err, &validationErr) {
			errs.BadRequest(w, r, err)
			return
		}

		errs.ServerError(w, r, err)
		return
	}

	rlog.FromContext(r.Context()).Info("updated movie", "movie", movie)

	err = body.WriteJSON(w, http.StatusOK, body.Envelope{"movie": movie}, nil)
	if err != nil {
		errs.ServerError(w, r, err)
		return
	}
}
