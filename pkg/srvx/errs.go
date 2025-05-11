package srvx

import (
	"fmt"
	"net/http"
)

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := Envelope{"error": message}
	Logger(r.Context()).Error("error response", "status", status, "message", message)

	err := WriteJSON(w, status, env, nil)
	if err != nil {
		LogErr(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ErrServer(w http.ResponseWriter, r *http.Request, err error) {
	LogErr(r, err)
	errorResponse(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func ErrNotFound(w http.ResponseWriter, r *http.Request) {
	errorResponse(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func ErrMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func ErrBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusBadRequest, err)
}
