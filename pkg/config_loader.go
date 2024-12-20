// config_loader.gp
package pkg

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the expected structure of config.yaml
type Config struct {
	OutputFormat string `yaml:"output_format"` // e.g., "json" or "markdown"
	ScanPath     string `yaml:"scan_path"`     // Directory to scan
	RulesPath    string `yaml:"rules_path"`    // Path to rules file
	ReportOutput string `yaml:"report_output"` // Path to save the report
	StrictMode   bool   `yaml:"strict_mode"`   // Enable strict validation
}

// DefaultConfig provides default values for config.yaml
func DefaultConfig() Config {
	return Config{
		OutputFormat: "json",
		ScanPath:     "./example-files",
		RulesPath:    "./custom-rules.yaml",
		ReportOutput: "./report.md",
		StrictMode:   false,
	}
}

// LoadConfig loads the configuration from a specified file or uses defaults.
func LoadConfig(configFile string) (Config, error) {
	// Start with default configuration
	config := DefaultConfig()

	// Attempt to load the config file if provided
	if configFile != "" {
		file, err := os.Open(configFile)
		if err != nil {
			log.Printf("Config file not found. Using defaults. Error: %v\n", err)
			return config, nil // Return defaults if file is missing
		}
		defer file.Close()

		decoder := yaml.NewDecoder(file)
		decoder.SetStrict(false) // Allow unknown fields without errors

		if err := decoder.Decode(&config); err != nil {
			return Config{}, fmt.Errorf("failed to parse config file: %v", err)
		}

		// log.Printf("Config file loaded successfully: %+v\n", config)
	}

	// Apply environment variable overrides
	if val, ok := os.LookupEnv("SCANRUNNER_OUTPUT_FORMAT"); ok {
		log.Printf("Overriding OutputFormat with environment variable: %s\n", val)
		config.OutputFormat = val
	}
	if val, ok := os.LookupEnv("SCANRUNNER_SCAN_PATH"); ok {
		log.Printf("Overriding ScanPath with environment variable: %s\n", val)
		config.ScanPath = val
	}
	if val, ok := os.LookupEnv("SCANRUNNER_RULES_PATH"); ok {
		log.Printf("Overriding RulesPath with environment variable: %s\n", val)
		config.RulesPath = val
	}
	if val, ok := os.LookupEnv("SCANRUNNER_REPORT_OUTPUT"); ok {
		log.Printf("Overriding ReportOutput with environment variable: %s\n", val)
		config.ReportOutput = val
	}
	if val, ok := os.LookupEnv("SCANRUNNER_STRICT_MODE"); ok {
		if val == "true" {
			log.Printf("Overriding StrictMode with environment variable: true\n")
			config.StrictMode = true
		} else if val == "false" {
			log.Printf("Overriding StrictMode with environment variable: false\n")
			config.StrictMode = false
		}
	}

	// log.Printf("Final Config: %+v\n", config)
	return config, nil
}
