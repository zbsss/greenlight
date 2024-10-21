package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/zbsss/greenlight/internal/model"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input model.CreateMovieParams

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err = input.OK(); err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	movie, err := app.db.CreateMovie(r.Context(), input)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.loggerFromContext(r).Info("created movie", "movie", movie)

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) listMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.db.ListMovies(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movies": movies}, nil)
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

	movie, err := app.db.GetMovie(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFound(w, r)
			return
		}

		app.serverError(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}
