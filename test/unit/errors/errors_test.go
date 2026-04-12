package errors_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apperrors "github.com/PlatformStackPulse/go-template/internal/errors"
)

func TestNew(t *testing.T) {
	err := apperrors.New(apperrors.ErrInvalidInput, "name is required")

	assert.NotNil(t, err)
	assert.Equal(t, apperrors.ErrInvalidInput, err.Code)
	assert.Equal(t, "name is required", err.Message)
	assert.Nil(t, err.Cause)
}

func TestError(t *testing.T) {
	tests := []struct {
		name     string
		err      *apperrors.AppError
		expected string
	}{
		{
			name:     "error without cause",
			err:      apperrors.New(apperrors.ErrInvalidInput, "invalid"),
			expected: "INVALID_INPUT: invalid",
		},
		{
			name:     "error with cause",
			err:      apperrors.Wrap(apperrors.ErrIntegration, "db failed", fmt.Errorf("connection refused")),
			expected: "INTEGRATION_ERROR: db failed (cause: connection refused)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestUnwrap(t *testing.T) {
	cause := fmt.Errorf("original error")
	err := apperrors.Wrap(apperrors.ErrIntegration, "integration failed", cause)

	assert.Equal(t, cause, err.Unwrap())
}

func TestWrap(t *testing.T) {
	cause := fmt.Errorf("database error")
	err := apperrors.Wrap(apperrors.ErrConfiguration, "config invalid", cause)

	assert.Equal(t, apperrors.ErrConfiguration, err.Code)
	assert.Equal(t, "config invalid", err.Message)
	assert.Equal(t, cause, err.Cause)
}

func TestIsCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		code     apperrors.Code
		expected bool
	}{
		{
			name:     "matching code",
			err:      apperrors.New(apperrors.ErrInvalidInput, "test"),
			code:     apperrors.ErrInvalidInput,
			expected: true,
		},
		{
			name:     "non-matching code",
			err:      apperrors.New(apperrors.ErrInvalidInput, "test"),
			code:     apperrors.ErrNotFound,
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			code:     apperrors.ErrInvalidInput,
			expected: false,
		},
		{
			name:     "wrapped error with matching code",
			err:      apperrors.Wrap(apperrors.ErrNotFound, "user not found", fmt.Errorf("sql error")),
			code:     apperrors.ErrNotFound,
			expected: true,
		},
		{
			name:     "non-AppError",
			err:      fmt.Errorf("generic error"),
			code:     apperrors.ErrInvalidInput,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := apperrors.IsCode(tt.err, tt.code)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestErrorChain(t *testing.T) {
	// Create nested errors
	sqlErr := fmt.Errorf("sql: connection refused")
	dbErr := apperrors.Wrap(apperrors.ErrIntegration, "database unavailable", sqlErr)
	appErr := apperrors.Wrap(apperrors.ErrInternal, "failed to save user", dbErr)

	// Check error chain
	assert.True(t, apperrors.IsCode(appErr, apperrors.ErrInternal))
	assert.True(t, apperrors.IsCode(appErr, apperrors.ErrIntegration), "should find wrapped errors")

	// Test errors.As
	var target *apperrors.AppError
	if errors.As(appErr, &target) {
		require.NotNil(t, target)
		assert.Equal(t, apperrors.ErrInternal, target.Code)
	} else {
		t.Fatal("failed to cast to AppError")
	}
}

func TestAllErrorCodes(t *testing.T) {
	// Ensure all error codes are defined and unique
	codes := []apperrors.Code{
		apperrors.ErrInvalidInput,
		apperrors.ErrNotFound,
		apperrors.ErrUnauthorized,
		apperrors.ErrConflict,
		apperrors.ErrInternal,
		apperrors.ErrConfiguration,
		apperrors.ErrIntegration,
		apperrors.ErrTimeout,
	}

	// Check uniqueness
	seen := make(map[apperrors.Code]bool)
	for _, code := range codes {
		require.False(t, seen[code], "duplicate error code: %v", code)
		seen[code] = true
	}
}
