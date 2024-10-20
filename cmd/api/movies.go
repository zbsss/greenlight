package main

import (
	"fmt"
	"net/http"

	"github.com/zbsss/greenlight/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Create movie")
}

func (app *application) viewMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	movie := data.Movie{
		ID:    id,
		Title: "Django",
		Genres: []string{
			"Tarantino",
			"Western",
		},
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}
