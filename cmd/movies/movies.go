package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zbsss/greenlight/internal/json"
	movies "github.com/zbsss/greenlight/internal/movies/service"
	"github.com/zbsss/greenlight/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var req movies.CreateMovieRequest

	err := json.Read(w, r, &req)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	movie, err := app.movies.CreateMovie(r.Context(), req)
	if err != nil {
		var validationErr validator.ValidationError
		if errors.As(err, &validationErr) {
			app.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		app.serverError(w, r, err)
		return
	}

	app.log.WithContext(r.Context()).Info("created movie", "movie", movie)

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = json.Write(w, http.StatusCreated, json.Envelope{"movie": movie}, headers)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) listMovies(w http.ResponseWriter, r *http.Request) {
	mvs, err := app.movies.ListMovies(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = json.Write(w, http.StatusOK, json.Envelope{"movies": mvs}, nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) viewMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		app.notFound(w, r)
		return
	}

	movie, err := app.movies.GetMovie(r.Context(), id)
	if err != nil {
		if errors.Is(err, movies.ErrMovieNotFound) {
			app.notFound(w, r)
			return
		}

		app.serverError(w, r, err)
		return
	}

	err = json.Write(w, http.StatusOK, json.Envelope{"movie": movie}, nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}
