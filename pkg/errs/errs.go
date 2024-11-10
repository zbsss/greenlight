package errs

import (
	"fmt"
	"net/http"

	"github.com/zbsss/greenlight/pkg/body"
	"github.com/zbsss/greenlight/pkg/rlog"
)

func logError(r *http.Request, err error) {
	rlog.FromContext(r.Context()).Error(err.Error(), "trace", fmt.Sprintf("%+v", err))
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := body.Envelope{"error": message}

	err := body.WriteJSON(w, status, env, nil)
	if err != nil {
		logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)
	errorResponse(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	errorResponse(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusBadRequest, err)
}
