package main

import (
	"context"
	"fmt"
	"os"

	"github.com/PlatformStackPulse/go-template/internal/cli"
	"github.com/PlatformStackPulse/go-template/internal/config"
	"github.com/PlatformStackPulse/go-template/internal/logger"
	"github.com/PlatformStackPulse/go-template/pkg/version"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.NewLogger(cfg.Debug)

	// Create root command
	cmd := cli.NewRootCommand(log)
	cmd.Version = version.Version

	// Add example command (remove or rename this for your project)
	cmd.AddCommand(cli.NewExampleCommand(log))

	// Execute command
	ctx := context.Background()
	if err := cmd.ExecuteContext(ctx); err != nil {
		log.Error("Command execution failed", "error", err)
		os.Exit(1) //nolint:gocritic
	}
}
