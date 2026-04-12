package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatformStackPulse/go-template/internal/domain"
)

func TestGreeterGreet(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "greet with name", input: "Alice", expected: "Hello, Alice!"},
		{name: "greet without name", input: "", expected: "Hello, World!"},
		{name: "greet with special characters", input: "Bob@123", expected: "Hello, Bob@123!"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			greeter := domain.NewGreeter(tc.input)
			assert.Equal(t, tc.expected, greeter.Greet())
		})
	}
}

func TestNewGreeter(t *testing.T) {
	greeter := domain.NewGreeter("TestName")
	assert.NotNil(t, greeter)
	assert.Equal(t, "TestName", greeter.Name)
}
