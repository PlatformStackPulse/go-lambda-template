// Package config provides configuration loading and management.
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds application configuration.
type Config struct {
	Debug                bool
	AppName              string
	Environment          string
	Version              string
	AWSRegion            string
	DynamoDBTableName    string
	GreetingParameterName string
}

// Load loads configuration from environment variables.
// Returns an error if required configuration is missing or invalid.
func Load() (*Config, error) {
	appName := getEnv("APP_NAME", "go-lambda-template")
	environment := getEnv("APP_ENV", "dev")
	greetingParameterName := getEnv("GREETING_PARAMETER_NAME", fmt.Sprintf("/%s/%s/greeting-prefix", appName, environment))

	cfg := &Config{
		Debug:                 getBoolEnv("DEBUG", false),
		AppName:               appName,
		Environment:           environment,
		Version:               getEnv("APP_VERSION", "dev"),
		AWSRegion:             getEnv("AWS_REGION", "us-east-1"),
		DynamoDBTableName:     getEnv("DYNAMODB_TABLE_NAME", fmt.Sprintf("%s-%s-requests", appName, environment)),
		GreetingParameterName: greetingParameterName,
	}

	// Validate required fields
	if cfg.AppName == "" {
		return nil, fmt.Errorf("APP_NAME cannot be empty")
	}
	if cfg.AWSRegion == "" {
		return nil, fmt.Errorf("AWS_REGION cannot be empty")
	}
	if cfg.DynamoDBTableName == "" {
		return nil, fmt.Errorf("DYNAMODB_TABLE_NAME cannot be empty")
	}
	if cfg.GreetingParameterName == "" {
		return nil, fmt.Errorf("GREETING_PARAMETER_NAME cannot be empty")
	}

	return cfg, nil
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
