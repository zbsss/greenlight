package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zbsss/greenlight/movies/backend/service"
	"github.com/zbsss/greenlight/pkg/srvx"
	"github.com/zbsss/greenlight/pkg/validator"
)

type Server struct {
	ms *service.MovieService
}

func NewServer(ms *service.MovieService) Server {
	return Server{ms: ms}
}

func (s Server) GetV1Movies(w http.ResponseWriter, r *http.Request) {
	mvs, err := s.ms.ListMovies(r.Context())
	if err != nil {
		srvx.ErrServer(w, r, err)
		return
	}

	apiMovies := make([]Movie, len(mvs))
	for i, m := range mvs {
		apiMovies[i] = *toAPIMovie(m)
	}

	if err := srvx.WriteJSON(w, http.StatusOK, srvx.Envelope{"movies": apiMovies}, nil); err != nil {
		srvx.ErrServer(w, r, err)
		return
	}
}

func (s Server) PostV1Movies(w http.ResponseWriter, r *http.Request) {
	var apiInput CreateMovieRequest
	err := srvx.ReadJSON(w, r, &apiInput)
	if err != nil {
		srvx.ErrBadRequest(w, r, err)
		return
	}

	movie, err := s.ms.CreateMovie(r.Context(), apiInput.toService())
	if err != nil {
		var validationErr validator.ValidationError
		if errors.As(err, &validationErr) {
			srvx.ErrBadRequest(w, r, err)
			return
		}

		srvx.ErrServer(w, r, err)
		return
	}

	srvx.Logger(r.Context()).Info("created movie", "movie", movie)

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	if err := srvx.WriteJSON(w, http.StatusCreated, srvx.Envelope{"movie": toAPIMovie(movie)}, headers); err != nil {
		srvx.ErrServer(w, r, err)
		return
	}
}

func (s Server) GetV1MoviesId(w http.ResponseWriter, r *http.Request, id int64) {
	movie, err := s.ms.GetMovie(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrMovieNotFound) {
			srvx.ErrNotFound(w, r)
			return
		}

		srvx.ErrServer(w, r, err)
		return
	}

	if err := srvx.WriteJSON(w, http.StatusOK, srvx.Envelope{"movie": toAPIMovie(movie)}, nil); err != nil {
		srvx.ErrServer(w, r, err)
		return
	}
}

func (s Server) PatchV1MoviesId(w http.ResponseWriter, r *http.Request, id int64) {
	var apiInput UpdateMovieRequest
	err := srvx.ReadJSON(w, r, &apiInput)
	if err != nil {
		srvx.ErrBadRequest(w, r, err)
		return
	}

	movie, err := s.ms.UpdateMovie(r.Context(), id, apiInput.toService())
	if err != nil {
		var validationErr validator.ValidationError
		if errors.As(err, &validationErr) {
			srvx.ErrBadRequest(w, r, err)
			return
		}

		srvx.ErrServer(w, r, err)
		return
	}

	srvx.Logger(r.Context()).Info("updated movie", "movie", movie)

	if err := srvx.WriteJSON(w, http.StatusOK, srvx.Envelope{"movie": movie}, nil); err != nil {
		srvx.ErrServer(w, r, err)
		return
	}
}
