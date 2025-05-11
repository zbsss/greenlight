package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zbsss/greenlight/movies/backend/storage"
	"github.com/zbsss/greenlight/movies/backend/storage/mocks"
	"github.com/zbsss/greenlight/pkg/validator"
	"k8s.io/utils/ptr"
)

var (
	errInjectedDBError = errors.New("injected database error")
)

// testHelpers contains common utilities for testing movie operations
type testHelpers struct {
	t       *testing.T
	model   *mocks.MockQueries
	service *MovieService
}

func setupTest(t *testing.T) testHelpers {
	mockModel := mocks.NewMockQueries()
	service := New(mockModel)
	return testHelpers{t: t, model: mockModel, service: service}
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

func TestCreateMovie(t *testing.T) {
	h := setupTest(t)

	validCreateMovie := MovieInput{
		Title:      "Casablanca",
		Year:       1942,
		RuntimeMin: 102,
		Genres:     []string{"drama", "romance", "war"},
	}
	expectedMovie := &Movie{
		ID:         1,
		Version:    1,
		Title:      "Casablanca",
		Year:       1942,
		RuntimeMin: 102,
		Genres:     []string{"drama", "romance", "war"},
	}

	tcs := []struct {
		name          string
		input         MovieInput
		injectDBError error
		expectedMovie *Movie
		expectedError error
	}{
		{
			name:          "valid input",
			input:         validCreateMovie,
			expectedMovie: expectedMovie,
		},
		{
			name:          "empty input",
			input:         MovieInput{},
			expectedError: validator.ValidationError{},
		},
		{
			name:          "db error",
			input:         validCreateMovie,
			injectDBError: errInjectedDBError,
			expectedError: errInjectedDBError,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(_ *testing.T) {
			h.model.Reset()
			if tc.injectDBError != nil {
				h.model.FailOnNextCall(tc.injectDBError)
			}

			actualMovie, err := h.service.CreateMovie(context.Background(), tc.input)
			h.assertError(tc.expectedError, err)
			h.assertMovie(tc.expectedMovie, actualMovie)
		})
	}
}
func TestUpdateMovie(t *testing.T) {
	h := setupTest(t)

	existingMovie := storage.Movie{
		ID:         1,
		Version:    1,
		Title:      "Django",
		Year:       2017,
		RuntimeMin: 120,
		Genres:     []string{"action"},
	}

	// Base expected movie that represents the movie with just the version incremented
	baseExpected := &Movie{
		ID:         1,
		Version:    2,
		Title:      "Django",
		Year:       2017,
		RuntimeMin: 120,
		Genres:     []string{"action"},
	}

	// Helper function to clone and modify the base expected movie
	cloneWithOverrides := func(overrides func(*Movie)) *Movie {
		clone := *baseExpected // Create a shallow copy
		if overrides != nil {
			overrides(&clone)
		}
		return &clone
	}

	tcs := []struct {
		name          string
		id            int64
		input         PartialMovieUpdate
		injectDBError error
		expectedMovie *Movie
		expectedError error
	}{
		{
			name: "update title",
			id:   1,
			input: PartialMovieUpdate{
				Title: ptr.To("Django Unchained"),
			},
			expectedMovie: cloneWithOverrides(func(m *Movie) {
				m.Title = "Django Unchained"
			}),
		},
		{
			name: "update year",
			id:   1,
			input: PartialMovieUpdate{
				Year: ptr.To[int32](2018),
			},
			expectedMovie: cloneWithOverrides(func(m *Movie) {
				m.Year = 2018
			}),
		},
		{
			name: "update runtime",
			id:   1,
			input: PartialMovieUpdate{
				RuntimeMin: ptr.To[int32](121),
			},
			expectedMovie: cloneWithOverrides(func(m *Movie) {
				m.RuntimeMin = 121
			}),
		},
		{
			name: "update genres",
			id:   1,
			input: PartialMovieUpdate{
				Genres: []string{"comedy"},
			},
			expectedMovie: cloneWithOverrides(func(m *Movie) {
				m.Genres = []string{"comedy"}
			}),
		},
		{
			name:          "empty update",
			id:            1,
			input:         PartialMovieUpdate{},
			expectedMovie: cloneWithOverrides(nil), // Use base expected without modifications
		},
		{
			name:          "not found",
			id:            2,
			input:         PartialMovieUpdate{},
			expectedError: ErrMovieNotFound,
		},
		{
			name:          "db error",
			input:         PartialMovieUpdate{},
			injectDBError: errInjectedDBError,
			expectedError: errInjectedDBError,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(_ *testing.T) {
			h.model.Reset(existingMovie)
			if tc.injectDBError != nil {
				h.model.FailOnNextCall(tc.injectDBError)
			}

			actualMovie, err := h.service.UpdateMovie(context.Background(), tc.id, tc.input)
			h.assertError(tc.expectedError, err)
			h.assertMovie(tc.expectedMovie, actualMovie)
		})
	}
}
