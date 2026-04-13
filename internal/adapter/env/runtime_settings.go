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

func (r *RuntimeSettings) Lookup(_ context.Context, key string) string {
	return strings.TrimSpace(os.Getenv(key))
}
