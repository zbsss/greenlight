package main

import (
	"net/http"

	"github.com/zbsss/greenlight/internal/json"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := json.Envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		}}

	err := json.Write(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}
