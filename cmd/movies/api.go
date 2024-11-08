package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/zbsss/greenlight/internal/errs"
	"github.com/zbsss/greenlight/internal/json"
	movies "github.com/zbsss/greenlight/internal/movies/service"
	"github.com/zbsss/greenlight/internal/rlog"
	"github.com/zbsss/greenlight/internal/validator"
)

func bindMoviesAPI(app *application, router *httprouter.Router) {
	api := moviesAPI{app: app}

	router.HandlerFunc("GET", "/v1/healthcheck", api.healthcheck)
	router.HandlerFunc("POST", "/v1/movies", api.create)
	router.HandlerFunc("GET", "/v1/movies", api.list)
	router.HandlerFunc("GET", "/v1/movies/:id", api.view)
}

type moviesAPI struct {
	app *application
}

func (api *moviesAPI) healthcheck(w http.ResponseWriter, r *http.Request) {
	data := json.Envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": api.app.config.env,
			"version":     version,
		}}

	err := json.Write(w, http.StatusOK, data, nil)
	if err != nil {
		errs.ServerError(w, r, err)
		return
	}
}

func (api *moviesAPI) create(w http.ResponseWriter, r *http.Request) {
	var req movies.CreateMovieRequest

	err := json.Read(w, r, &req)
	if err != nil {
		errs.BadRequest(w, r, err)
		return
	}

	movie, err := api.app.movies.CreateMovie(r.Context(), req)
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

	err = json.Write(w, http.StatusCreated, json.Envelope{"movie": movie}, headers)
	if err != nil {
		errs.ServerError(w, r, err)
		return
	}
}

func (api *moviesAPI) list(w http.ResponseWriter, r *http.Request) {
	mvs, err := api.app.movies.ListMovies(r.Context())
	if err != nil {
		errs.ServerError(w, r, err)
		return
	}

	err = json.Write(w, http.StatusOK, json.Envelope{"movies": mvs}, nil)
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

	movie, err := api.app.movies.GetMovie(r.Context(), id)
	if err != nil {
		if errors.Is(err, movies.ErrMovieNotFound) {
			errs.NotFound(w, r)
			return
		}

		errs.ServerError(w, r, err)
		return
	}

	err = json.Write(w, http.StatusOK, json.Envelope{"movie": movie}, nil)
	if err != nil {
		errs.ServerError(w, r, err)
		return
	}
}
