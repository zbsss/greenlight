package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	movies "github.com/zbsss/greenlight/internal/movies/service"
	"github.com/zbsss/greenlight/pkg/body"
	"github.com/zbsss/greenlight/pkg/errs"
	"github.com/zbsss/greenlight/pkg/rlog"
	"github.com/zbsss/greenlight/pkg/validator"
)

func BindMoviesAPI(ms *movies.MovieService, router *httprouter.Router) {
	api := moviesAPI{ms: ms}

	router.HandlerFunc("POST", "/v1/movies", api.create)
	router.HandlerFunc("GET", "/v1/movies", api.list)
	router.HandlerFunc("GET", "/v1/movies/:id", api.view)
	router.HandlerFunc("PATCH", "/v1/movies/:id", api.update)
}

type moviesAPI struct {
	ms *movies.MovieService
}

func (api *moviesAPI) create(w http.ResponseWriter, r *http.Request) {
	var input movies.CreateMovieRequest
	err := body.ReadJSON(w, r, &input)
	if err != nil {
		errs.BadRequest(w, r, err)
		return
	}

	movie, err := api.ms.CreateMovie(r.Context(), input)
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

func (api *moviesAPI) list(w http.ResponseWriter, r *http.Request) {
	mvs, err := api.ms.ListMovies(r.Context())
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

func (api *moviesAPI) view(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		errs.NotFound(w, r)
		return
	}

	movie, err := api.ms.GetMovie(r.Context(), id)
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

func (api *moviesAPI) update(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		errs.NotFound(w, r)
		return
	}

	var input movies.UpdateMovieRequest
	err = body.ReadJSON(w, r, &input)
	if err != nil {
		errs.BadRequest(w, r, err)
		return
	}

	movie, err := api.ms.UpdateMovie(r.Context(), id, input)
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
