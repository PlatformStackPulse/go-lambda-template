package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/PlatformStackPulse/go-lambda-template/internal/domain"
	apperrors "github.com/PlatformStackPulse/go-lambda-template/internal/errors"
	"github.com/PlatformStackPulse/go-lambda-template/internal/logger"
	"github.com/PlatformStackPulse/go-lambda-template/internal/usecase"
)

type stubConfigProvider struct {
	prefix string
	err    error
}

func (s stubConfigProvider) StringValue(context.Context, string) (string, error) {
	return s.prefix, s.err
}

type stubRecorder struct {
	record domain.GreetingRecord
	err    error
}

func (s *stubRecorder) Record(_ context.Context, record domain.GreetingRecord) error {
	s.record = record
	return s.err
}

type stubEnvironment struct {
	prefixOverride string
	sourceLabel    string
}

func (s stubEnvironment) Lookup(_ context.Context, key string) string {
	switch key {
	case "SAMPLE_GREETING_PREFIX":
		return s.prefixOverride
	case "API_SOURCE_LABEL":
		return s.sourceLabel
	default:
		return ""
	}
}

func TestGreetingUseCaseExecute(t *testing.T) {
	// Main flow should resolve prefix, persist a record, and return a response DTO.
	log := logger.NewLogger(false)
	recorder := &stubRecorder{}
	now := func() time.Time {
		return time.Date(2026, time.April, 12, 10, 30, 0, 0, time.UTC)
	}
	uc := usecase.NewGreetingUseCaseWithClock(log, stubConfigProvider{prefix: "Hello"}, recorder, stubEnvironment{}, now)

	result, err := uc.Execute(context.Background(), usecase.GreetingInput{
		Name:      "Charlie",
		RequestID: "req-123",
		Source:    "unit-test",
	})
	require.NoError(t, err)

	assert.Equal(t, "Hello, Charlie!", result.Message)
	assert.Equal(t, "req-123", result.RequestID)
	assert.Equal(t, "unit-test", result.Source)
	assert.Equal(t, "2026-04-12T10:30:00Z", result.Timestamp)
	assert.Equal(t, "Charlie", recorder.record.Name)
	assert.Equal(t, "Hello, Charlie!", recorder.record.Message)
}

func TestGreetingUseCaseFallsBackToWorld(t *testing.T) {
	// Empty name should use domain fallback behavior.
	log := logger.NewLogger(false)
	recorder := &stubRecorder{}
	uc := usecase.NewGreetingUseCase(log, stubConfigProvider{prefix: "Welcome"}, recorder, stubEnvironment{})

	result, err := uc.Execute(context.Background(), usecase.GreetingInput{})
	require.NoError(t, err)
	assert.Equal(t, "Welcome, World!", result.Message)
	assert.Equal(t, "unknown", result.RequestID)
	assert.Equal(t, "api", result.Source)
}

func TestGreetingUseCaseReturnsIntegrationErrorOnConfigFailure(t *testing.T) {
	log := logger.NewLogger(false)
	uc := usecase.NewGreetingUseCase(log, stubConfigProvider{err: errors.New("ssm unavailable")}, &stubRecorder{}, stubEnvironment{})

	_, err := uc.Execute(context.Background(), usecase.GreetingInput{Name: "Taylor"})
	require.Error(t, err)
	assert.True(t, apperrors.IsCode(err, apperrors.ErrIntegration))
}

func TestGreetingUseCaseReturnsIntegrationErrorOnRecorderFailure(t *testing.T) {
	log := logger.NewLogger(false)
	uc := usecase.NewGreetingUseCase(log, stubConfigProvider{prefix: "Hello"}, &stubRecorder{err: errors.New("ddb unavailable")}, stubEnvironment{})

	_, err := uc.Execute(context.Background(), usecase.GreetingInput{Name: "Taylor"})
	require.Error(t, err)
	assert.True(t, apperrors.IsCode(err, apperrors.ErrIntegration))
}

func TestGreetingUseCaseUsesEnvironmentOverride(t *testing.T) {
	// Environment overrides should take precedence over loaded configuration.
	log := logger.NewLogger(false)
	recorder := &stubRecorder{}
	uc := usecase.NewGreetingUseCase(log, stubConfigProvider{prefix: "Ignored"}, recorder, stubEnvironment{prefixOverride: "Howdy", sourceLabel: "env-source"})

	result, err := uc.Execute(context.Background(), usecase.GreetingInput{Name: "Taylor", RequestID: "req-env", Source: "request-source"})
	require.NoError(t, err)
	assert.Equal(t, "Howdy, Taylor!", result.Message)
	assert.Equal(t, "env-source", result.Source)
}
