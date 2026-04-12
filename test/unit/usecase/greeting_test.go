package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatformStackPulse/go-template/internal/logger"
	"github.com/PlatformStackPulse/go-template/internal/usecase"
)

func TestGreetingUseCaseExecute(t *testing.T) {
	log := logger.NewLogger(false)
	uc := usecase.NewGreetingUseCase(log)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "greeting with name", input: "Charlie", expected: "Hello, Charlie!"},
		{name: "greeting without name", input: "", expected: "Hello, World!"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := uc.Execute(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}
