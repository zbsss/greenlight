package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/zbsss/greenlight/pkg/body"
	"github.com/zbsss/greenlight/pkg/errs"
)

func bindHealthAPI(app *application, router *httprouter.Router) {
	api := healthAPI{app: app}

	router.HandlerFunc("GET", "/v1/healthcheck", api.healthcheck)
}

type healthAPI struct {
	app *application
}

func (api *healthAPI) healthcheck(w http.ResponseWriter, r *http.Request) {
	data := body.Envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": api.app.config.env,
			"version":     version,
		}}

	err := body.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		errs.ServerError(w, r, err)
		return
	}
}
