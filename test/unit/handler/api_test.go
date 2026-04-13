package handler_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apperrors "github.com/PlatformStackPulse/go-lambda-template/internal/errors"
	"github.com/PlatformStackPulse/go-lambda-template/internal/handler"
	"github.com/PlatformStackPulse/go-lambda-template/internal/logger"
	"github.com/PlatformStackPulse/go-lambda-template/internal/usecase"
)

type stubExecutor struct {
	output usecase.GreetingOutput
	err    error
	input  usecase.GreetingInput
}

func (s *stubExecutor) Execute(_ context.Context, input usecase.GreetingInput) (usecase.GreetingOutput, error) {
	s.input = input
	return s.output, s.err
}

func TestAPIHandlerSuccess(t *testing.T) {
	// Successful execution should map to HTTP 200 with JSON payload.
	executor := &stubExecutor{output: usecase.GreetingOutput{Message: "Hello, Sam!", RequestID: "req-1", Source: "$default", Timestamp: "2026-04-12T10:00:00Z"}}
	h := handler.NewAPIHandler(logger.NewLogger(false), executor)

	response, err := h.Handle(context.Background(), events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{"name": "Sam"},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			RequestID: "req-1",
			Stage:     "$default",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "Sam", executor.input.Name)

	var payload usecase.GreetingOutput
	require.NoError(t, json.Unmarshal([]byte(response.Body), &payload))
	assert.Equal(t, "Hello, Sam!", payload.Message)
}

func TestAPIHandlerReturnsClientMessageForKnownErrors(t *testing.T) {
	// Known business errors should expose user-safe messages.
	executor := &stubExecutor{err: apperrors.New(apperrors.ErrInvalidInput, "name is invalid")}
	h := handler.NewAPIHandler(logger.NewLogger(false), executor)

	response, err := h.Handle(context.Background(), events.APIGatewayV2HTTPRequest{})
	require.NoError(t, err)
	assert.Equal(t, 400, response.StatusCode)
	assert.Contains(t, response.Body, "name is invalid")
}

func TestAPIHandlerReturnsGenericMessageForInternalErrors(t *testing.T) {
	// Internal/integration errors should hide internals behind a generic response.
	executor := &stubExecutor{err: apperrors.New(apperrors.ErrIntegration, "ssm unavailable")}
	h := handler.NewAPIHandler(logger.NewLogger(false), executor)

	response, err := h.Handle(context.Background(), events.APIGatewayV2HTTPRequest{})
	require.NoError(t, err)
	assert.Equal(t, 500, response.StatusCode)
	assert.Contains(t, response.Body, "internal server error")
}

func TestAPIHandlerHealthRoute(t *testing.T) {
	executor := &stubExecutor{}
	h := handler.NewAPIHandler(logger.NewLogger(false), executor)

	response, err := h.Handle(context.Background(), events.APIGatewayV2HTTPRequest{
		RawPath: "/health",
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			RequestID: "req-health",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Contains(t, response.Body, "\"check\":\"liveness\"")
	assert.Equal(t, "", executor.input.Name)
}

func TestAPIHandlerReadyRoute(t *testing.T) {
	executor := &stubExecutor{}
	h := handler.NewAPIHandler(logger.NewLogger(false), executor)

	response, err := h.Handle(context.Background(), events.APIGatewayV2HTTPRequest{
		RawPath: "/ready",
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			RequestID: "req-ready",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Contains(t, response.Body, "\"check\":\"readiness\"")
	assert.Equal(t, "", executor.input.Name)
}
