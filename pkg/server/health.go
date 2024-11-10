package server

import (
	"net/http"

	"github.com/zbsss/greenlight/pkg/body"
	"github.com/zbsss/greenlight/pkg/errs"
)

func healthcheck(w http.ResponseWriter, r *http.Request) {
	data := body.Envelope{
		"status": "OK",
	}

	err := body.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		errs.ServerError(w, r, err)
		return
	}
}
