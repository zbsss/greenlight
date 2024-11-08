package main

import (
	"fmt"
	"net/http"

	"github.com/zbsss/greenlight/internal/json"
)

func (app *application) logError(r *http.Request, err error) {
	app.log.WithContext(r.Context()).Error(err.Error(), "trace", fmt.Sprintf("%+v", err))
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := json.Envelope{"error": message}

	err := json.Write(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	app.errorResponse(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}
