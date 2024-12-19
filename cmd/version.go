// version.go - Implementation for the "version" command in SCANRUNNER-CLI

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Define the version of the CLI tool
const version = "v1.0.0"

// versionCmd represents the "version" subcommand
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the current version of SCANRUNNER-CLI",
	Long:  `The version command outputs the current version of the SCANRUNNER-CLI tool.`,
	// The Run function should contain logic to:
	// - Print the version of the tool to the terminal.
	Run: func(cmd *cobra.Command, args []string) {
		// Display the version to the user.
		fmt.Printf("SCANRUNNER-CLI version: %s\n", version)
	},
}

func init() {
	// Register the "version" command as a subcommand of the root.
	rootCmd.AddCommand(versionCmd)
}
