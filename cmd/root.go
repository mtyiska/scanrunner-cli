// root.go - Base setup for the SCANRUNNER-CLI tool

package cmd

import (
	"fmt"
	"os"

	"github.com/mtyiska/scanrunner/pkg"
	"github.com/spf13/cobra"
)

var (
	configFile string    // Variable to hold the path to the config file
	config     pkg.Config // Variable to store the loaded configuration
	rules      pkg.Rules  // Variable to store the loaded rules
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scanrunner",
	Short: "SCANRUNNER-CLI: A modular tool for file validation and AI-powered insights",
	Long: `SCANRUNNER-CLI is a flexible command-line tool for scanning, validating,
and reporting on YAML/JSON files while leveraging AI for actionable suggestions.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load the configuration file
		var err error
		config, err = pkg.LoadConfig(configFile)
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// Load the rules file using the path specified in the config
		rules, err = pkg.LoadRules(config.RulesPath)
		if err != nil {
			return fmt.Errorf("failed to load rules: %w", err)
		}

		// log.Printf("Configuration loaded: %+v\n", config)
		// log.Printf("Rules loaded: %+v\n", rules)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is the entry point for the CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Add a persistent flag for specifying the configuration file
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "./config/default-config.yaml",
    	"Path to the configuration file (default is ./config/default-config.yaml)")

}
