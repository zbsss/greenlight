package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/zbsss/greenlight/internal/movies/model/mocks"
	movies "github.com/zbsss/greenlight/internal/movies/service"
	"github.com/zbsss/greenlight/internal/server/testserver"
)

func TestGetMovie(t *testing.T) {
	mockDB := mocks.NewMockQueries()
	ms := movies.NewMovieService(mockDB)
	router := httprouter.New()

	BindMoviesAPI(ms, router)

	ts := testserver.New(t, router)
	defer ts.Close()

	tcs := []struct {
		name           string
		id             int
		injectDBError  error
		expectedStatus int
	}{
		{
			name:           "ok",
			id:             1,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "not found",
			id:             2,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "db error",
			id:             1,
			injectDBError:  fmt.Errorf("something went wrong"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			mockDB.Reset()
			if tc.injectDBError != nil {
				mockDB.FailOnNextCall(tc.injectDBError)
			}

			url := fmt.Sprintf("/v1/movies/%d", tc.id)
			code, headers, _ := ts.Get(t, url)

			if code != tc.expectedStatus {
				t.Fatalf("expected status %d, got %d", tc.expectedStatus, code)
			}

			if headers.Get("Content-Type") != "application/json" {
				t.Error("expected Content-Type header to be application/json")
			}
		})
	}
}
