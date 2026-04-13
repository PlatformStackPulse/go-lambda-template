package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"

	apperrors "github.com/PlatformStackPulse/go-lambda-template/internal/errors"
	"github.com/PlatformStackPulse/go-lambda-template/internal/logger"
	"github.com/PlatformStackPulse/go-lambda-template/internal/usecase"
	"github.com/PlatformStackPulse/go-lambda-template/pkg/health"
)

type GreetingExecutor interface {
	Execute(context.Context, usecase.GreetingInput) (usecase.GreetingOutput, error)
}

type APIHandler struct {
	log      *logger.Logger
	executor GreetingExecutor
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewAPIHandler(log *logger.Logger, executor GreetingExecutor) *APIHandler {
	return &APIHandler{log: log, executor: executor}
}

func (h *APIHandler) Handle(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	requestPath := extractRequestPath(request)
	if requestPath == "/health" {
		return jsonResponse(http.StatusOK, health.Health(time.Now()), request.RequestContext.RequestID), nil
	}
	if requestPath == "/ready" {
		return jsonResponse(http.StatusOK, health.Ready(time.Now()), request.RequestContext.RequestID), nil
	}

	// Support both path and query input so local and deployed calls behave consistently.
	name := request.PathParameters["name"]
	if name == "" {
		name = request.QueryStringParameters["name"]
	}

	output, err := h.executor.Execute(ctx, usecase.GreetingInput{
		Name:      name,
		RequestID: request.RequestContext.RequestID,
		Source:    request.RequestContext.Stage,
	})
	if err != nil {
		// Domain/application errors are translated into HTTP-safe status and messages.
		h.log.Error("request failed", "error", err)
		return jsonResponse(mapErrorStatus(err), errorResponse{Error: userMessage(err)}, request.RequestContext.RequestID), nil
	}

	return jsonResponse(http.StatusOK, output, request.RequestContext.RequestID), nil
}

func extractRequestPath(request events.APIGatewayV2HTTPRequest) string {
	path := request.RawPath
	if path == "" {
		path = request.RequestContext.HTTP.Path
	}
	if path == "" {
		return ""
	}

	normalized := strings.TrimRight(path, "/")
	if normalized == "" {
		return "/"
	}

	return normalized
}

func jsonResponse(status int, payload any, requestID string) events.APIGatewayV2HTTPResponse {
	body, err := json.Marshal(payload)
	if err != nil {
		// Always return valid JSON even when response marshaling unexpectedly fails.
		body = []byte(`{"error":"failed to encode response"}`)
		status = http.StatusInternalServerError
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Request-Id": requestID,
		},
		Body: string(body),
	}
}

func mapErrorStatus(err error) int {
	switch {
	case apperrors.IsCode(err, apperrors.ErrInvalidInput):
		return http.StatusBadRequest
	case apperrors.IsCode(err, apperrors.ErrUnauthorized):
		return http.StatusUnauthorized
	case apperrors.IsCode(err, apperrors.ErrNotFound):
		return http.StatusNotFound
	case apperrors.IsCode(err, apperrors.ErrConflict):
		return http.StatusConflict
	case apperrors.IsCode(err, apperrors.ErrTimeout):
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}

func userMessage(err error) string {
	// Only expose explicit user-safe messages for expected business errors.
	if ok := apperrors.IsCode(err, apperrors.ErrInvalidInput) ||
		apperrors.IsCode(err, apperrors.ErrNotFound) ||
		apperrors.IsCode(err, apperrors.ErrConflict) ||
		apperrors.IsCode(err, apperrors.ErrUnauthorized); ok {
		var typedErr *apperrors.AppError
		if errors.As(err, &typedErr) && typedErr.Message != "" {
			return typedErr.Message
		}
	}

	return "internal server error"
}
