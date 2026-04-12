// Package config provides configuration loading and management.
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds application configuration.
type Config struct {
	Debug   bool
	AppName string
	Version string
}

// Load loads configuration from environment variables.
// Returns an error if required configuration is missing or invalid.
func Load() (*Config, error) {
	cfg := &Config{
		Debug:   getBoolEnv("DEBUG", false),
		AppName: getEnv("APP_NAME", "go-template"),
		Version: getEnv("APP_VERSION", "dev"),
	}

	// Validate required fields
	if cfg.AppName == "" {
		return nil, fmt.Errorf("APP_NAME cannot be empty")
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
