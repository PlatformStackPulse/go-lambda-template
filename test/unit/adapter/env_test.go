package adapter_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	envadapter "github.com/PlatformStackPulse/go-lambda-template/internal/adapter/env"
)

func TestRuntimeSettingsReadsEnvironmentVariables(t *testing.T) {
	_ = os.Setenv("SAMPLE_GREETING_PREFIX", "Howdy")
	_ = os.Setenv("API_SOURCE_LABEL", "lambda-env")
	defer func() {
		_ = os.Unsetenv("SAMPLE_GREETING_PREFIX")
		_ = os.Unsetenv("API_SOURCE_LABEL")
	}()

	settings := envadapter.NewRuntimeSettings()
	assert.Equal(t, "Howdy", settings.Lookup(context.Background(), "SAMPLE_GREETING_PREFIX"))
	assert.Equal(t, "lambda-env", settings.Lookup(context.Background(), "API_SOURCE_LABEL"))
}

func TestRuntimeSettingsTrimsAndDefaults(t *testing.T) {
	_ = os.Setenv("SAMPLE_GREETING_PREFIX", " ")
	_ = os.Setenv("API_SOURCE_LABEL", "  api-source ")
	defer func() {
		_ = os.Unsetenv("SAMPLE_GREETING_PREFIX")
		_ = os.Unsetenv("API_SOURCE_LABEL")
	}()

	settings := envadapter.NewRuntimeSettings()
	assert.Equal(t, "", settings.Lookup(context.Background(), "SAMPLE_GREETING_PREFIX"))
	assert.Equal(t, "api-source", settings.Lookup(context.Background(), "API_SOURCE_LABEL"))
}
