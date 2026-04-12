// Package cli provides command-line interface commands and initialization.
package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/PlatformStackPulse/go-template/internal/logger"
)

// NewExampleCommand creates an example command.
// This demonstrates how to add commands to your CLI.
// Remove or rename this command for your application.
func NewExampleCommand(log *logger.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hello",
		Short: "Example command",
		Long:  `An example command demonstrating command structure. Remove this and add your own commands.`,
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")

			if name != "" {
				fmt.Printf("Hello, %s!\n", name)
				log.Info("Hello executed", "name", name)
			} else {
				fmt.Println("Hello, World!")
				log.Info("Hello executed")
			}
		},
	}

	cmd.Flags().StringP("name", "n", "", "Name to greet")

	return cmd
}
