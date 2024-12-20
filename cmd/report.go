package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mtyiska/scanrunner/internal/compliance"
	"github.com/mtyiska/scanrunner/pkg"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a summary report from validation results",
	Long: `The report command aggregates and formats validation results from the
	validate command into a readable report. Reports can be output as JSON, Markdown,
	or other formats, based on user preference.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Generating report in format: %s\n", config.OutputFormat)
		log.Printf("Saving report to: %s\n", config.ReportOutput)

		// Scan the directory for YAML files to validate
		files, err := scanDirectory(config.ScanPath)
		if err != nil {
			log.Fatalf("Error scanning directory: %v\n", err)
		}
		if len(files) == 0 {
			log.Println("No files found for validation. Report generation skipped.")
			return
		}

		// Load validation rules
		rules, err := pkg.LoadRules(config.RulesPath)
		if err != nil {
			log.Fatalf("Failed to load validation rules: %v\n", err)
		}

		// Validate each file and store results
		var results []map[string]interface{}
		for _, file := range files {
			// log.Printf("Validating file: %s\n", file)
			err := compliance.ValidateFile(file, rules)
			if err != nil {
				// log.Printf("Validation failed for file %s: %v\n", file, err)
				results = append(results, map[string]interface{}{
					"file":   filepath.Base(file),
					"status": "FAIL",
					"error":  err.Error(),
				})
				continue
			}

			results = append(results, map[string]interface{}{
				"file":   filepath.Base(file),
				"status": "PASS",
			})
		}

		// Format and save the report
		reportContent := formatReport(results, config.OutputFormat)
		err = saveReport(reportContent, config.ReportOutput)
		if err != nil {
			log.Fatalf("Error saving report: %v\n", err)
		}

		fmt.Println("Report successfully generated.")
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}

// formatReport formats the validation results based on the desired output format
func formatReport(results []map[string]interface{}, format string) string {
	switch format {
	case "json":
		jsonContent, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			log.Printf("Error formatting JSON report: %v\n", err)
			return ""
		}
		return string(jsonContent)
	case "markdown":
		report := "# Validation Report\n\n"
		for _, result := range results {
			file := result["file"].(string)
			status := result["status"].(string)
			report += fmt.Sprintf("- **%s**: %s\n", file, status)
			if status == "FAIL" {
				errorMsg := result["error"].(string)
				report += fmt.Sprintf("  - %s\n", errorMsg)
			}
		}
		return report
	default:
		log.Printf("Unknown format: %s. Defaulting to markdown.\n", format)
		return formatReport(results, "markdown")
	}
}

// saveReport writes the report content to a specified file
func saveReport(content, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}
