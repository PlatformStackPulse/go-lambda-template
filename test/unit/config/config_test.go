package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/PlatformStackPulse/go-lambda-template/internal/config"
)

func TestLoadDefaults(t *testing.T) {
	_ = os.Unsetenv("DEBUG")
	_ = os.Unsetenv("APP_NAME")
	_ = os.Unsetenv("APP_ENV")
	_ = os.Unsetenv("APP_VERSION")
	_ = os.Unsetenv("AWS_REGION")
	_ = os.Unsetenv("DYNAMODB_TABLE_NAME")
	_ = os.Unsetenv("GREETING_PARAMETER_NAME")

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.False(t, cfg.Debug)
	assert.Equal(t, "go-lambda-template", cfg.AppName)
	assert.Equal(t, "dev", cfg.Environment)
	assert.Equal(t, "dev", cfg.Version)
	assert.Equal(t, "us-east-1", cfg.AWSRegion)
	assert.Equal(t, "go-lambda-template-dev-requests", cfg.DynamoDBTableName)
	assert.Equal(t, "/go-lambda-template/dev/greeting-prefix", cfg.GreetingParameterName)
}

func TestLoadFromEnv(t *testing.T) {
	_ = os.Setenv("DEBUG", "true")
	_ = os.Setenv("APP_NAME", "my-app")
	_ = os.Setenv("APP_ENV", "prod")
	_ = os.Setenv("APP_VERSION", "v1.0.0")
	_ = os.Setenv("AWS_REGION", "eu-west-1")
	_ = os.Setenv("DYNAMODB_TABLE_NAME", "sample-table")
	_ = os.Setenv("GREETING_PARAMETER_NAME", "/sample/prod/greeting-prefix")
	defer func() {
		_ = os.Unsetenv("DEBUG")
		_ = os.Unsetenv("APP_NAME")
		_ = os.Unsetenv("APP_ENV")
		_ = os.Unsetenv("APP_VERSION")
		_ = os.Unsetenv("AWS_REGION")
		_ = os.Unsetenv("DYNAMODB_TABLE_NAME")
		_ = os.Unsetenv("GREETING_PARAMETER_NAME")
	}()

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.True(t, cfg.Debug)
	assert.Equal(t, "my-app", cfg.AppName)
	assert.Equal(t, "prod", cfg.Environment)
	assert.Equal(t, "v1.0.0", cfg.Version)
	assert.Equal(t, "eu-west-1", cfg.AWSRegion)
	assert.Equal(t, "sample-table", cfg.DynamoDBTableName)
	assert.Equal(t, "/sample/prod/greeting-prefix", cfg.GreetingParameterName)
}

func TestInvalidBoolEnvFallsBack(t *testing.T) {
	_ = os.Setenv("DEBUG", "notabool")
	defer func() { _ = os.Unsetenv("DEBUG") }()

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.False(t, cfg.Debug)
}
