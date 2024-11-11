package movies

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zbsss/greenlight/internal/movies/model/mocks"
	"github.com/zbsss/greenlight/pkg/validator"
)

var (
	errInjectedDBError = errors.New("injected database error")
)

// testCase represents a common test case structure for movie service operations
type testCase struct {
	name          string
	id            int64
	input         MovieInput
	injectDBError error
	expectedMovie *Movie
	expectedError error
}

// testHelpers contains common utilities for testing movie operations
type testHelpers struct {
	t       *testing.T
	model   *mocks.MockQueries
	service *MovieService
}

func setupTest(t *testing.T) testHelpers {
	mockModel := mocks.NewMockQueries()
	service := NewMovieService(mockModel)
	return testHelpers{t: t, model: mockModel, service: service}
}

func (h testHelpers) runTestCase(tc *testCase, operation func(context.Context, int64, MovieInput) (*Movie, error)) {
	h.t.Helper()

	h.model.Reset()
	if tc.injectDBError != nil {
		h.model.FailOnNextCall(tc.injectDBError)
	}

	actualMovie, err := operation(context.Background(), tc.id, tc.input)
	h.assertError(tc.expectedError, err)
	h.assertMovie(tc.expectedMovie, actualMovie)
}

func (h testHelpers) assertError(expected, actual error) {
	h.t.Helper()

	if expected != nil {
		if actual == nil {
			h.t.Fatalf("expected error to be %v; got nil", expected)
		}

		if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
			h.t.Fatalf("expected error to be %T; got %T", expected, actual)
		}
	} else if actual != nil {
		h.t.Fatalf("did not expect an error but got %v", actual)
	}
}

func (h testHelpers) assertMovie(expected, actual *Movie) {
	h.t.Helper()

	if !cmp.Equal(actual, expected) {
		h.t.Fatalf("expected movie to be %v; got %v", expected, actual)
	}
}

// Test data factories
func validMovieInput() MovieInput {
	return MovieInput{
		Title:      "Casablanca",
		Year:       1942,
		RuntimeMin: 102,
		Genres:     []string{"drama", "romance", "war"},
	}
}

func movieResult(id int64, version int32) *Movie {
	return &Movie{
		ID:      id,
		Version: version,
		Title:   "Casablanca",
		Year:    1942,
		Runtime: Runtime(102),
		Genres:  []string{"drama", "romance", "war"},
	}
}

func TestCreateMovie(t *testing.T) {
	h := setupTest(t)

	tcs := []testCase{
		{
			name:          "valid input",
			input:         validMovieInput(),
			expectedMovie: movieResult(2, 1),
		},
		{
			name:          "empty input",
			input:         MovieInput{},
			expectedError: validator.ValidationError{},
		},
		{
			name:          "db error",
			input:         validMovieInput(),
			injectDBError: errInjectedDBError,
			expectedError: errInjectedDBError,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(_ *testing.T) {
			h.runTestCase(&tc, func(ctx context.Context, _ int64, input MovieInput) (*Movie, error) {
				return h.service.CreateMovie(ctx, input)
			})
		})
	}
}

func TestUpdateMovie(t *testing.T) {
	h := setupTest(t)

	tcs := []testCase{
		{
			name:          "valid input",
			id:            1,
			input:         validMovieInput(),
			expectedMovie: movieResult(1, 2),
		},
		{
			name:          "empty input",
			id:            1,
			input:         MovieInput{},
			expectedError: validator.ValidationError{},
		},
		{
			name:          "not found",
			id:            2,
			input:         validMovieInput(),
			expectedError: ErrMovieNotFound,
		},
		{
			name:          "db error",
			input:         validMovieInput(),
			injectDBError: errInjectedDBError,
			expectedError: errInjectedDBError,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(_ *testing.T) {
			h.runTestCase(&tc, h.service.UpdateMovie)
		})
	}
}
