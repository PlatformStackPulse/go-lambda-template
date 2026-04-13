package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/PlatformStackPulse/go-lambda-template/internal/config"
)

func TestLoadDefaults(t *testing.T) {
	// Clear all supported variables to verify deterministic defaults.
	_ = os.Unsetenv("DEBUG")
	_ = os.Unsetenv("APP_NAME")
	_ = os.Unsetenv("APP_ENV")
	_ = os.Unsetenv("APP_VERSION")
	_ = os.Unsetenv("AWS_REGION")
	_ = os.Unsetenv("API_BASE_PATH")
	_ = os.Unsetenv("API_SOURCE_LABEL")
	_ = os.Unsetenv("DYNAMODB_REQUESTS_TABLE_NAME")
	_ = os.Unsetenv("SSM_PARAMETER_PREFIX")
	_ = os.Unsetenv("APP_CONFIG_PARAMETER_NAME")
	_ = os.Unsetenv("S3_SOURCE_BUCKET_NAME")
	_ = os.Unsetenv("S3_SOURCE_KEY_PREFIX")
	_ = os.Unsetenv("S3_TARGET_BUCKET_NAME")
	_ = os.Unsetenv("S3_TARGET_KEY_PREFIX")
	_ = os.Unsetenv("KMS_KEY_ARN")
	_ = os.Unsetenv("POSTGRES_ENABLED")
	_ = os.Unsetenv("POSTGRES_DATABASE_NAME")
	_ = os.Unsetenv("POSTGRES_CLUSTER_ARN")
	_ = os.Unsetenv("POSTGRES_SECRET_ARN")

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.False(t, cfg.Debug)
	assert.Equal(t, "go-lambda-template", cfg.App.Name)
	assert.Equal(t, "dev", cfg.App.Environment)
	assert.Equal(t, "dev", cfg.App.Version)
	assert.Equal(t, "us-east-1", cfg.AWS.Region)
	assert.Equal(t, "/hello", cfg.API.BasePath)
	assert.Equal(t, "", cfg.API.SourceLabel)
	assert.Equal(t, "go-lambda-template-dev-requests", cfg.DynamoDB.RequestsTableName)
	assert.Equal(t, "/go-lambda-template/dev", cfg.SSM.ParameterPrefix)
	assert.Equal(t, "/go-lambda-template/dev/app-config", cfg.SSM.AppConfigParameterName)
	assert.False(t, cfg.Postgres.Enabled)
	assert.Equal(t, "app", cfg.Postgres.DatabaseName)
	assert.Equal(t, "", cfg.Postgres.ClusterARN)
	assert.Equal(t, "", cfg.Postgres.SecretARN)
}

func TestLoadFromEnv(t *testing.T) {
	// Set every supported override to ensure env mapping is complete.
	_ = os.Setenv("DEBUG", "true")
	_ = os.Setenv("APP_NAME", "my-app")
	_ = os.Setenv("APP_ENV", "prod")
	_ = os.Setenv("APP_VERSION", "v1.0.0")
	_ = os.Setenv("AWS_REGION", "eu-west-1")
	_ = os.Setenv("API_BASE_PATH", "/requests")
	_ = os.Setenv("API_SOURCE_LABEL", "edge")
	_ = os.Setenv("DYNAMODB_REQUESTS_TABLE_NAME", "sample-table")
	_ = os.Setenv("SSM_PARAMETER_PREFIX", "/sample/prod")
	_ = os.Setenv("APP_CONFIG_PARAMETER_NAME", "/sample/prod/platform-config")
	_ = os.Setenv("S3_SOURCE_BUCKET_NAME", "source-bucket")
	_ = os.Setenv("S3_SOURCE_KEY_PREFIX", "incoming/")
	_ = os.Setenv("S3_TARGET_BUCKET_NAME", "target-bucket")
	_ = os.Setenv("S3_TARGET_KEY_PREFIX", "processed/")
	_ = os.Setenv("KMS_KEY_ARN", "arn:aws:kms:eu-west-1:123456789012:key/example")
	_ = os.Setenv("POSTGRES_ENABLED", "true")
	_ = os.Setenv("POSTGRES_DATABASE_NAME", "orders")
	_ = os.Setenv("POSTGRES_CLUSTER_ARN", "arn:aws:rds:eu-west-1:123456789012:cluster:example")
	_ = os.Setenv("POSTGRES_SECRET_ARN", "arn:aws:secretsmanager:eu-west-1:123456789012:secret:example")
	defer func() {
		_ = os.Unsetenv("DEBUG")
		_ = os.Unsetenv("APP_NAME")
		_ = os.Unsetenv("APP_ENV")
		_ = os.Unsetenv("APP_VERSION")
		_ = os.Unsetenv("AWS_REGION")
		_ = os.Unsetenv("API_BASE_PATH")
		_ = os.Unsetenv("API_SOURCE_LABEL")
		_ = os.Unsetenv("DYNAMODB_REQUESTS_TABLE_NAME")
		_ = os.Unsetenv("SSM_PARAMETER_PREFIX")
		_ = os.Unsetenv("APP_CONFIG_PARAMETER_NAME")
		_ = os.Unsetenv("S3_SOURCE_BUCKET_NAME")
		_ = os.Unsetenv("S3_SOURCE_KEY_PREFIX")
		_ = os.Unsetenv("S3_TARGET_BUCKET_NAME")
		_ = os.Unsetenv("S3_TARGET_KEY_PREFIX")
		_ = os.Unsetenv("KMS_KEY_ARN")
		_ = os.Unsetenv("POSTGRES_ENABLED")
		_ = os.Unsetenv("POSTGRES_DATABASE_NAME")
		_ = os.Unsetenv("POSTGRES_CLUSTER_ARN")
		_ = os.Unsetenv("POSTGRES_SECRET_ARN")
	}()

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.True(t, cfg.Debug)
	assert.Equal(t, "my-app", cfg.App.Name)
	assert.Equal(t, "prod", cfg.App.Environment)
	assert.Equal(t, "v1.0.0", cfg.App.Version)
	assert.Equal(t, "eu-west-1", cfg.AWS.Region)
	assert.Equal(t, "/requests", cfg.API.BasePath)
	assert.Equal(t, "edge", cfg.API.SourceLabel)
	assert.Equal(t, "sample-table", cfg.DynamoDB.RequestsTableName)
	assert.Equal(t, "/sample/prod", cfg.SSM.ParameterPrefix)
	assert.Equal(t, "/sample/prod/platform-config", cfg.SSM.AppConfigParameterName)
	assert.Equal(t, "source-bucket", cfg.S3.SourceBucketName)
	assert.Equal(t, "incoming/", cfg.S3.SourceKeyPrefix)
	assert.Equal(t, "target-bucket", cfg.S3.TargetBucketName)
	assert.Equal(t, "processed/", cfg.S3.TargetKeyPrefix)
	assert.Equal(t, "arn:aws:kms:eu-west-1:123456789012:key/example", cfg.KMS.KeyARN)
	assert.True(t, cfg.Postgres.Enabled)
	assert.Equal(t, "orders", cfg.Postgres.DatabaseName)
	assert.Equal(t, "arn:aws:rds:eu-west-1:123456789012:cluster:example", cfg.Postgres.ClusterARN)
	assert.Equal(t, "arn:aws:secretsmanager:eu-west-1:123456789012:secret:example", cfg.Postgres.SecretARN)
}

func TestInvalidBoolEnvFallsBack(t *testing.T) {
	_ = os.Setenv("DEBUG", "notabool")
	defer func() { _ = os.Unsetenv("DEBUG") }()

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.False(t, cfg.Debug)
}

func TestPostgresValidationRequiresArnsWhenEnabled(t *testing.T) {
	// Postgres mode must require both cluster and secret ARN wiring.
	_ = os.Setenv("POSTGRES_ENABLED", "true")
	defer func() {
		_ = os.Unsetenv("POSTGRES_ENABLED")
		_ = os.Unsetenv("POSTGRES_CLUSTER_ARN")
		_ = os.Unsetenv("POSTGRES_SECRET_ARN")
		_ = os.Unsetenv("POSTGRES_DATABASE_NAME")
	}()

	_, err := config.Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "POSTGRES_CLUSTER_ARN")
}

func TestDefaultAppConfigDocument(t *testing.T) {
	cfg, err := config.Load()
	require.NoError(t, err)

	document := config.DefaultAppConfigDocument(cfg)
	assert.Equal(t, "/hello", document["api.base_path"])
	assert.Equal(t, "go-lambda-template", document["platform.name"])
	assert.Equal(t, "Hello", document["sample.greeting.prefix"])
}
