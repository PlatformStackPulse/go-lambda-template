// Package config provides configuration loading and management.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds application configuration.
type Config struct {
	Debug    bool
	App      AppConfig
	AWS      AWSConfig
	API      APIConfig
	DynamoDB DynamoDBConfig
	SSM      SSMConfig
	S3       S3Config
	KMS      KMSConfig
	Postgres PostgresConfig
}

type AppConfig struct {
	Name        string
	Environment string
	Version     string
}

type AWSConfig struct {
	Region string
}

type APIConfig struct {
	BasePath    string
	SourceLabel string
}

type DynamoDBConfig struct {
	RequestsTableName string
}

type SSMConfig struct {
	ParameterPrefix        string
	AppConfigParameterName string
}

type S3Config struct {
	SourceBucketName string
	SourceKeyPrefix  string
	TargetBucketName string
	TargetKeyPrefix  string
}

type KMSConfig struct {
	KeyARN string
}

type PostgresConfig struct {
	Enabled      bool
	DatabaseName string
	ClusterARN   string
	SecretARN    string
}

// Load loads configuration from environment variables.
// Returns an error if required configuration is missing or invalid.
func Load() (*Config, error) {
	appName := getEnv("APP_NAME", "go-lambda-template")
	environment := getEnv("APP_ENV", "dev")
	parameterPrefix := getEnv("SSM_PARAMETER_PREFIX", fmt.Sprintf("/%s/%s", appName, environment))
	appConfigParameterName := getEnv("APP_CONFIG_PARAMETER_NAME", fmt.Sprintf("%s/app-config", strings.TrimRight(parameterPrefix, "/")))

	cfg := &Config{
		Debug: getBoolEnv("DEBUG", false),
		App: AppConfig{
			Name:        appName,
			Environment: environment,
			Version:     getEnv("APP_VERSION", "dev"),
		},
		AWS: AWSConfig{
			Region: getEnv("AWS_REGION", "us-east-1"),
		},
		API: APIConfig{
			BasePath:    getEnv("API_BASE_PATH", "/hello"),
			SourceLabel: getEnv("API_SOURCE_LABEL", ""),
		},
		DynamoDB: DynamoDBConfig{
			RequestsTableName: getEnv("DYNAMODB_REQUESTS_TABLE_NAME", fmt.Sprintf("%s-%s-requests", appName, environment)),
		},
		SSM: SSMConfig{
			ParameterPrefix:        parameterPrefix,
			AppConfigParameterName: appConfigParameterName,
		},
		S3: S3Config{
			SourceBucketName: getEnv("S3_SOURCE_BUCKET_NAME", ""),
			SourceKeyPrefix:  getEnv("S3_SOURCE_KEY_PREFIX", ""),
			TargetBucketName: getEnv("S3_TARGET_BUCKET_NAME", ""),
			TargetKeyPrefix:  getEnv("S3_TARGET_KEY_PREFIX", ""),
		},
		KMS: KMSConfig{
			KeyARN: getEnv("KMS_KEY_ARN", ""),
		},
		Postgres: PostgresConfig{
			Enabled:      getBoolEnv("POSTGRES_ENABLED", false),
			DatabaseName: getEnv("POSTGRES_DATABASE_NAME", "app"),
			ClusterARN:   getEnv("POSTGRES_CLUSTER_ARN", ""),
			SecretARN:    getEnv("POSTGRES_SECRET_ARN", ""),
		},
	}

	// Validate required fields
	if cfg.App.Name == "" {
		return nil, fmt.Errorf("APP_NAME cannot be empty")
	}
	if cfg.AWS.Region == "" {
		return nil, fmt.Errorf("AWS_REGION cannot be empty")
	}
	if cfg.DynamoDB.RequestsTableName == "" {
		return nil, fmt.Errorf("DYNAMODB_REQUESTS_TABLE_NAME cannot be empty")
	}
	if cfg.SSM.ParameterPrefix == "" {
		return nil, fmt.Errorf("SSM_PARAMETER_PREFIX cannot be empty")
	}
	if cfg.SSM.AppConfigParameterName == "" {
		return nil, fmt.Errorf("APP_CONFIG_PARAMETER_NAME cannot be empty")
	}
	if cfg.API.BasePath == "" {
		return nil, fmt.Errorf("API_BASE_PATH cannot be empty")
	}
	if cfg.Postgres.Enabled {
		if cfg.Postgres.DatabaseName == "" {
			return nil, fmt.Errorf("POSTGRES_DATABASE_NAME cannot be empty when POSTGRES_ENABLED=true")
		}
		if cfg.Postgres.ClusterARN == "" {
			return nil, fmt.Errorf("POSTGRES_CLUSTER_ARN cannot be empty when POSTGRES_ENABLED=true")
		}
		if cfg.Postgres.SecretARN == "" {
			return nil, fmt.Errorf("POSTGRES_SECRET_ARN cannot be empty when POSTGRES_ENABLED=true")
		}
	}

	if err := validateAppConfigJSON(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateAppConfigJSON(cfg *Config) error {
	defaultDocument := DefaultAppConfigDocument(cfg)
	if _, err := json.Marshal(defaultDocument); err != nil {
		return fmt.Errorf("default app config document is invalid: %w", err)
	}

	return nil
}

func DefaultAppConfigDocument(cfg *Config) map[string]string {
	return map[string]string{
		"api.base_path":                cfg.API.BasePath,
		"dynamodb.requests_table_name": cfg.DynamoDB.RequestsTableName,
		"platform.environment":         cfg.App.Environment,
		"platform.name":                cfg.App.Name,
		"platform.version":             cfg.App.Version,
		"sample.greeting.prefix":       "Hello",
		"s3.source.bucket":             cfg.S3.SourceBucketName,
		"s3.source.key_prefix":         cfg.S3.SourceKeyPrefix,
		"s3.target.bucket":             cfg.S3.TargetBucketName,
		"s3.target.key_prefix":         cfg.S3.TargetKeyPrefix,
		"kms.key_arn":                  cfg.KMS.KeyARN,
	}
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// getBoolEnv retrieves a boolean environment variable or returns a default value.
func getBoolEnv(key string, defaultVal bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultVal
}
