package main

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}
