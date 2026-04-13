package domain

import (
	"fmt"
	"strings"
)

type GreetingRecord struct {
	// RequestID is the traceable identifier propagated from API Gateway.
	RequestID string `dynamodbav:"request_id" json:"request_id"`
	// Name is the normalized caller name used in the response.
	Name string `dynamodbav:"name" json:"name"`
	// Message is the rendered greeting text returned to the client.
	Message string `dynamodbav:"message" json:"message"`
	// CreatedAt stores the UTC RFC3339 timestamp captured by the use case.
	CreatedAt string `dynamodbav:"created_at" json:"created_at"`
	// Source identifies the API stage or override source label.
	Source string `dynamodbav:"source" json:"source"`
}

func NormalizeName(name string) string {
	// Treat blank input as an anonymous request to keep responses stable.
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "World"
	}

	return trimmed
}

func BuildGreeting(prefix, name string) string {
	// Fall back to a safe default when configuration does not provide a prefix.
	trimmedPrefix := strings.TrimSpace(prefix)
	if trimmedPrefix == "" {
		trimmedPrefix = "Hello"
	}

	return fmt.Sprintf("%s, %s!", trimmedPrefix, NormalizeName(name))
}
