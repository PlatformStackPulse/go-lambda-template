package env

import (
	"context"
	"os"
	"strings"
)

type RuntimeSettings struct{}

func NewRuntimeSettings() *RuntimeSettings {
	return &RuntimeSettings{}
}

func (r *RuntimeSettings) GreetingPrefixOverride(_ context.Context) string {
	return strings.TrimSpace(os.Getenv("GREETING_PREFIX_OVERRIDE"))
}

func (r *RuntimeSettings) GreetingSourceLabel(_ context.Context) string {
	return strings.TrimSpace(os.Getenv("GREETING_SOURCE_LABEL"))
}
