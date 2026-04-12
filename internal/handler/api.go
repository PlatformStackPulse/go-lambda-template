package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	apperrors "github.com/PlatformStackPulse/go-lambda-template/internal/errors"
	"github.com/PlatformStackPulse/go-lambda-template/internal/logger"
	"github.com/PlatformStackPulse/go-lambda-template/internal/usecase"
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
		h.log.Error("request failed", "error", err)
		return jsonResponse(mapErrorStatus(err), errorResponse{Error: userMessage(err)}, request.RequestContext.RequestID), nil
	}

	return jsonResponse(http.StatusOK, output, request.RequestContext.RequestID), nil
}

func jsonResponse(status int, payload any, requestID string) events.APIGatewayV2HTTPResponse {
	body, err := json.Marshal(payload)
	if err != nil {
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
