// Package cli provides command-line interface commands and initialization.
package cli

import (
	"github.com/spf13/cobra"

	"github.com/PlatformStackPulse/go-template/internal/logger"
)

// NewRootCommand creates the root command
func NewRootCommand(log *logger.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go-template",
		Short: "Go Template CLI",
		Long:  `A production-ready Go CLI template for building scalable applications.`,
	}

	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		log.Debug("Command started", "command", cmd.Name())
	}

	return cmd
}
