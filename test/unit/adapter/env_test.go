package adapter_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	envadapter "github.com/PlatformStackPulse/go-lambda-template/internal/adapter/env"
)

func TestRuntimeSettingsReadsEnvironmentVariables(t *testing.T) {
	_ = os.Setenv("GREETING_PREFIX_OVERRIDE", "Howdy")
	_ = os.Setenv("GREETING_SOURCE_LABEL", "lambda-env")
	defer func() {
		_ = os.Unsetenv("GREETING_PREFIX_OVERRIDE")
		_ = os.Unsetenv("GREETING_SOURCE_LABEL")
	}()

	settings := envadapter.NewRuntimeSettings()
	assert.Equal(t, "Howdy", settings.GreetingPrefixOverride(context.Background()))
	assert.Equal(t, "lambda-env", settings.GreetingSourceLabel(context.Background()))
}

func TestRuntimeSettingsTrimsAndDefaults(t *testing.T) {
	_ = os.Setenv("GREETING_PREFIX_OVERRIDE", " ")
	_ = os.Setenv("GREETING_SOURCE_LABEL", "  api-source ")
	defer func() {
		_ = os.Unsetenv("GREETING_PREFIX_OVERRIDE")
		_ = os.Unsetenv("GREETING_SOURCE_LABEL")
	}()

	settings := envadapter.NewRuntimeSettings()
	assert.Equal(t, "", settings.GreetingPrefixOverride(context.Background()))
	assert.Equal(t, "api-source", settings.GreetingSourceLabel(context.Background()))
}
