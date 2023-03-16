package validate

import (
	"encoding/json"
	"errors"
)

// ErrInvalidId occurs when an ID is not in a valid form.
var ErrInvalidId = errors.New("ID is not in its proper form")

// ErrorResponse is the form used for API responses from failures in the API.
type ErrorResponse struct {
	Error  string `json:"error"`
	Fields string `json:"fields,omitempty"`
}

// RequestError is used to pass an error during the request through the app with
// web specific context.
type RequestError struct {
	Err    error
	Status int
	Fields error
}

// NewRequestError wraps a provided error with an HTTP status code.
func NewRequestError(err error, status int) error {
	return &RequestError{err, status, nil}
}

// Error implements the error interface.
func (e *RequestError) Error() string {
	return e.Err.Error()
}

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// FieldErrors represents a collection of field errors.
type FieldErrors []FieldError

// Error implements the error interface.
func (e FieldErrors) Error() string {
	d, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}

	return string(d)
}

// Cause iterates through all the wrapped errors until the root error value is
// reached.
func Cause(err error) error {
	root := err
	for {
		if err = errors.Unwrap(root); err == nil {
			return root
		}
		root = err
	}
}
