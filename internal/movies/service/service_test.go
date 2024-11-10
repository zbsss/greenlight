package movies

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zbsss/greenlight/internal/movies/model/mocks"
	"github.com/zbsss/greenlight/internal/validator"
)

func TestCreateMovie(t *testing.T) {
	mockModel := mocks.NewMockQueries()
	service := NewMovieService(mockModel)

	tcs := []struct {
		name          string
		input         CreateMovieRequest
		injectDBError error
		expectedMovie *Movie
		expectedError error
	}{
		{
			name: "valid input",
			input: CreateMovieRequest{
				Title:      "Casablanca",
				Year:       1942,
				RuntimeMin: 102,
				Genres:     []string{"drama", "romance", "war"},
			},
			expectedMovie: &Movie{
				ID:      2,
				Title:   "Casablanca",
				Year:    1942,
				Runtime: 102,
				Genres:  []string{"drama", "romance", "war"},
			},
		},
		{
			name:          "empty input",
			input:         CreateMovieRequest{},
			expectedError: validator.ValidationError{},
		},
		{
			name: "db error",
			input: CreateMovieRequest{
				Title:      "Casablanca",
				Year:       1942,
				RuntimeMin: 102,
				Genres:     []string{"drama", "romance", "war"},
			},
			injectDBError: errors.New("something went wrong"),
			expectedError: errors.New("something went wrong"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			mockModel.Reset()
			if tc.injectDBError != nil {
				mockModel.FailOnNextCall(tc.injectDBError)
			}

			actualMovie, err := service.CreateMovie(context.Background(), tc.input)
			if tc.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error to be %v; got nil", tc.expectedError)
				}

				if reflect.TypeOf(err) != reflect.TypeOf(tc.expectedError) {
					t.Fatalf("expected error to be %T; got %T", tc.expectedError, err)
				}
			} else if err != nil {
				t.Fatalf("did not expect an error but got %v", err)
			}

			if !cmp.Equal(actualMovie, tc.expectedMovie) {
				t.Fatalf("expected movie to be %v; got %v", tc.expectedMovie, actualMovie)
			}
		})
	}
}
