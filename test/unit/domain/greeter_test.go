package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatformStackPulse/go-lambda-template/internal/domain"
)

func TestNormalizeName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "trim name", input: " Alice ", expected: "Alice"},
		{name: "fallback name", input: "", expected: "World"},
		{name: "keep special characters", input: "Bob@123", expected: "Bob@123"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, domain.NormalizeName(tc.input))
		})
	}
}

func TestBuildGreeting(t *testing.T) {
	assert.Equal(t, "Hello, Alice!", domain.BuildGreeting("Hello", "Alice"))
	assert.Equal(t, "Hello, World!", domain.BuildGreeting("", ""))
}
