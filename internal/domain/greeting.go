package domain

import (
	"fmt"
	"strings"
)

type GreetingRecord struct {
	RequestID string `dynamodbav:"request_id" json:"request_id"`
	Name      string `dynamodbav:"name" json:"name"`
	Message   string `dynamodbav:"message" json:"message"`
	CreatedAt string `dynamodbav:"created_at" json:"created_at"`
	Source    string `dynamodbav:"source" json:"source"`
}

func NormalizeName(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "World"
	}

	return trimmed
}

func BuildGreeting(prefix, name string) string {
	trimmedPrefix := strings.TrimSpace(prefix)
	if trimmedPrefix == "" {
		trimmedPrefix = "Hello"
	}

	return fmt.Sprintf("%s, %s!", trimmedPrefix, NormalizeName(name))
}