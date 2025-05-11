package srvx

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testStruct struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Valid bool   `json:"valid"`
}

func TestReadJSON(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		wantErr        bool
		wantErrMessage string
		wantStruct     testStruct
	}{
		{
			name:       "valid JSON",
			body:       `{"name": "John", "age": 30, "valid": true}`,
			wantErr:    false,
			wantStruct: testStruct{Name: "John", Age: 30, Valid: true},
		},
		{
			name:           "invalid JSON syntax",
			body:           `{"name": "John", "age": 30, "valid": true`,
			wantErr:        true,
			wantErrMessage: "body contains badly-formed JSON",
		},
		{
			name:           "empty body",
			body:           "",
			wantErr:        true,
			wantErrMessage: "body must not be empty",
		},
		{
			name:           "unknown field",
			body:           `{"name": "John", "age": 30, "valid": true, "unknown": "field"}`,
			wantErr:        true,
			wantErrMessage: "body contains unknown key \"unknown\"",
		},
		{
			name:           "incorrect JSON type",
			body:           `{"name": "John", "age": "not a number", "valid": true}`,
			wantErr:        true,
			wantErrMessage: "body contains incorrect JSON type for field \"age\"",
		},
		{
			name:           "body size limit exceeded",
			body:           string(bytes.Repeat([]byte{'a'}, 1048577)),
			wantErr:        true,
			wantErrMessage: "body contains badly-formed JSON (at character 1)",
		},
		{
			name:           "multiple JSON values",
			body:           `{"name": "John", "age": 30, "valid": true}{"name": "Jane", "age": 25, "valid": true}`,
			wantErr:        true,
			wantErrMessage: "body must only contain a single JSON value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()
			var got testStruct

			err := ReadJSON(w, req, &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err.Error() != tt.wantErrMessage {
					t.Errorf("ReadJSON() error message = %v, want %v", err.Error(), tt.wantErrMessage)
				}
			} else {
				if got != tt.wantStruct {
					t.Errorf("ReadJSON() got = %v, want %v", got, tt.wantStruct)
				}
			}
		})
	}
}
