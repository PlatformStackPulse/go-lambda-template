package lambda_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/PlatformStackPulse/go-lambda-template/internal/handler"
	"github.com/PlatformStackPulse/go-lambda-template/internal/logger"
	"github.com/PlatformStackPulse/go-lambda-template/internal/usecase"
)

type integrationExecutor struct{}

func (integrationExecutor) Execute(_ context.Context, input usecase.GreetingInput) (usecase.GreetingOutput, error) {
	return usecase.GreetingOutput{
		Message:   "Hello, " + input.Name + "!",
		RequestID: input.RequestID,
		Source:    input.Source,
		Timestamp: "2026-04-12T12:00:00Z",
	}, nil
}

func TestAPIHandlerWithFixture(t *testing.T) {
	// Exercise the handler with a realistic API Gateway payload fixture.
	bytes, err := os.ReadFile(filepath.Join("..", "..", "fixtures", "events", "apigw-request.json"))
	require.NoError(t, err)

	var request events.APIGatewayV2HTTPRequest
	require.NoError(t, json.Unmarshal(bytes, &request))

	h := handler.NewAPIHandler(logger.NewLogger(false), integrationExecutor{})
	response, err := h.Handle(context.Background(), request)
	require.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Contains(t, response.Body, "Taylor")
	assert.Equal(t, "fixture-request-id", response.Headers["X-Request-Id"])
}
