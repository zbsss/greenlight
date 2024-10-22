package validator

type ValidationError struct {
	Errors map[string]string `json:"fieldErrors"`
}

// Implement the error interface by providing the Error() method
func (ve ValidationError) Error() string {
	return "validation errors"
}
