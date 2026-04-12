// Package errors provides application-specific error types and utilities.
package errors

import (
	"fmt"
)

// Code represents an application error code for categorization.
type Code string

const (
	// ErrInvalidInput indicates invalid user input.
	ErrInvalidInput Code = "INVALID_INPUT"

	// ErrNotFound indicates a resource was not found.
	ErrNotFound Code = "NOT_FOUND"

	// ErrUnauthorized indicates unauthorized access.
	ErrUnauthorized Code = "UNAUTHORIZED"

	// ErrConflict indicates a conflict (e.g., resource already exists).
	ErrConflict Code = "CONFLICT"

	// ErrInternal indicates an internal application error.
	ErrInternal Code = "INTERNAL_ERROR"

	// ErrConfiguration indicates a configuration problem.
	ErrConfiguration Code = "CONFIG_ERROR"

	// ErrIntegration indicates an error integrating with external service.
	ErrIntegration Code = "INTEGRATION_ERROR"

	// ErrTimeout indicates an operation timed out.
	ErrTimeout Code = "TIMEOUT"
)

// AppError is a structured application error.
type AppError struct {
	Code    Code   // Error code for categorization
	Message string // User-facing error message
	Cause   error  // Underlying error (can be nil)
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error for error chain analysis.
func (e *AppError) Unwrap() error {
	return e.Cause
}

// New creates a new AppError with the given code and message.
func New(code Code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   nil,
	}
}

// Wrap wraps an existing error with an application error code and message.
func Wrap(code Code, message string, cause error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// IsCode checks if an error matches a specific error code.
// It recursively unwraps errors to find a matching AppError code.
func IsCode(err error, code Code) bool {
	if err == nil {
		return false
	}

	// Check if this error is directly an AppError with matching code
	if appErr, ok := err.(*AppError); ok {
		if appErr.Code == code {
			return true
		}
		// Also check the wrapped error
		if appErr.Cause != nil {
			return IsCode(appErr.Cause, code)
		}
		return false
	}

	// Try to unwrap and check the underlying error
	if unwrappable, ok := err.(interface{ Unwrap() error }); ok {
		return IsCode(unwrappable.Unwrap(), code)
	}

	return false
}
