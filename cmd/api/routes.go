package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/zbsss/greenlight/pkg/middleware"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.NotFound = http.HandlerFunc(app.notFound)

	router.HandlerFunc("GET", "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc("POST", "/v1/movies", app.createMovieHandler)
	router.HandlerFunc("GET", "/v1/movies", app.listMovies)
	router.HandlerFunc("GET", "/v1/movies/:id", app.viewMovieHandler)

	// standard middleware for all requests
	standard := alice.New(app.recoverPanic, middleware.TraceMiddleware, app.logRequest, secureHeaders)

	return standard.Then(router)
}
