package web

import "errors"

// shutdownError is a type used to help with the graceful service termination.
type shutdownError struct {
	Message string
}

// NewShutdownError returns an error that causes a shutdown signal.
func NewShutdownError(msg string) error {
	return &shutdownError{Message: msg}
}

// Error is the implementation of the error interface.
func (s *shutdownError) Error() string {
	return s.Message
}

// IsShutdown checks to see if the shutdown error is contained in the error.
func IsShutdown(err error) bool {
	var s *shutdownError
	return errors.As(err, &s)
}
