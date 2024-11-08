package validator

import (
	"regexp"
	"slices"
)

type Valid interface {
	OK() error
}

// Define a new Validator type which contains a map of validation errors
type Validator struct {
	errors map[string]string
}

// New is a helper which creates a new Validator instance with an empty errors map.
func New() *Validator {
	return &Validator{errors: make(map[string]string)}
}

func (v *Validator) OK() error {
	if len(v.errors) == 0 {
		return nil
	}

	return ValidationError{Errors: v.errors}
}

// AddError adds an error message to the map (so long as no entry already exists for
// the given key).
func (v *Validator) AddError(key, message string) {
	if _, exists := v.errors[key]; !exists {
		v.errors[key] = message
	}
}

// Check adds an error message to the map only if a validation check is not 'ok'.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Generic function which returns true if a specific value is in a list of permitted
// values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Matches returns true if a string value matches a specific regexp pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Generic function which returns true if all values in a slice are unique.
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}