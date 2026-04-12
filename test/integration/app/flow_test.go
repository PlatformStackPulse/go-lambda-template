package app_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatformStackPulse/go-template/internal/logger"
	"github.com/PlatformStackPulse/go-template/internal/usecase"
)

func TestGreetingFlowIntegration(t *testing.T) {
	log := logger.NewLogger(false)
	uc := usecase.NewGreetingUseCase(log)

	msg, err := uc.Execute("Platform")
	assert.NoError(t, err)
	assert.Equal(t, "Hello, Platform!", msg)
}
