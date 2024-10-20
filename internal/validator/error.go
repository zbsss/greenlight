package validator

import "encoding/json"

type ValidationError struct {
	errors map[string]string
}

// Implement the error interface by providing the Error() method
func (ve ValidationError) Error() string {
	return "validation errors"
}

// Implement the json.Marshaler interface for custom JSON encoding
func (ve ValidationError) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(ve.errors, "", "  ")
}
