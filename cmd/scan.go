package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan directories for YAML/JSON files",
	Long: `The scan command identifies YAML and JSON files in the directory specified
	by the configuration file or via the --path flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Scanning directory: %s\n", config.ScanPath)

		// Scan the directory specified in config.ScanPath
		files, err := scanDirectory(config.ScanPath)
		if err != nil {
			log.Fatalf("Error scanning directory: %v\n", err)
		}

		// Output the discovered files
		if len(files) == 0 {
			fmt.Println("No YAML/JSON files found.")
		} else {
			fmt.Println("Discovered files:")
			for _, file := range files {
				fmt.Printf("- %s\n", file)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}

// scanDirectory scans the directory for YAML and JSON files
func scanDirectory(path string) ([]string, error) {
	var files []string

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (filepath.Ext(filePath) == ".yaml" || filepath.Ext(filePath) == ".yml") {
			files = append(files, filePath)
		}
		return nil
	})

	return files, err
}
