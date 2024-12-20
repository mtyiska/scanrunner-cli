// Validate.go
package cmd

import (
	"fmt"
	"log"

	"github.com/mtyiska/scanrunner/internal/compliance"
	"github.com/mtyiska/scanrunner/pkg"

	"github.com/spf13/cobra"
)

var strictMode bool
var rulesPath string

// validateCmd represents the "validate" subcommand
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate YAML, JSON files, and Dockerfiles against compliance rules",
	Long: `The validate command checks YAML, JSON files, and Dockerfiles for structural correctness
	and compliance with predefined rules (e.g., field presence, best practices, etc.).`,
	Run: func(cmd *cobra.Command, args []string) {
		// log.Printf("Validating files in directory: %s\n", config.ScanPath)

		// Load compliance rules
		if rulesPath == "" {
			rulesPath = config.RulesPath
		}
		rules, err := pkg.LoadRules(rulesPath)
		if err != nil {
			log.Fatalf("Failed to load rules: %v\n", err)
		}
		// log.Printf("Loaded rules: %+v\n", rules)

		// Scan the directory for YAML files
		files, err := scanDirectory(config.ScanPath)
		if err != nil {
			log.Fatalf("Error scanning directory: %v\n", err)
		}
		if len(files) == 0 {
			log.Println("No YAML, JSON files, or Dockerfiles found for validation.")
			return
		}

		// Validate each file against the rules
		var validationErrors []string
		for _, file := range files {
			// log.Printf("Validating file: %s\n", file)
			err := compliance.ValidateFile(file, rules)
			if err != nil {
				// log.Printf("Validation failed for %s: %v\n", file, err)
				validationErrors = append(validationErrors, fmt.Sprintf("%s: FAIL (%v)", file, err))
				if strictMode {
					log.Fatal("Strict mode enabled. Stopping on first error.")
				}
			} else {
				// log.Printf("Validation passed for %s\n", file)
				validationErrors = append(validationErrors, fmt.Sprintf("%s: PASS", file))
			}
		}

		// Print final validation results
		fmt.Println("\nValidation Results:")
		for _, result := range validationErrors {
			fmt.Println(result)
		}
	},
}

func init() {
	// Register flags
	validateCmd.Flags().BoolVar(&strictMode, "strict", false, "Enable strict mode for validation")
	validateCmd.Flags().StringVarP(&rulesPath, "rules", "r", "", "Path to custom compliance rules file")

	// Register the validate command
	rootCmd.AddCommand(validateCmd)
}
