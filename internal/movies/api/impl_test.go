package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/zbsss/greenlight/internal/movies/model/mocks"
	movies "github.com/zbsss/greenlight/internal/movies/service"
	"github.com/zbsss/greenlight/pkg/srvx/testserver"
)

func TestGetMovie(t *testing.T) {
	db := mocks.NewMockQueries()
	ms := movies.NewMovieService(db)
	router := http.NewServeMux()
	movieServer := NewServer(ms)
	h := HandlerFromMux(movieServer, router)

	ts := testserver.New(h)
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
			db.Reset(mocks.TestMovie1)
			if tc.injectDBError != nil {
				db.FailOnNextCall(tc.injectDBError)
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
