// ValidateFile.go
package pkg

import (
	"fmt"
)

// ValidateFile validates a single YAML file against the given rules
func ValidateFile(filePath string, rules Rules) error {
	// Parse and convert the YAML file
	parsedData, err := ParseAndConvertYAML(filePath)
	if err != nil {
		return fmt.Errorf("error parsing YAML file: %w", err)
	}

	// Validate each required field
	for _, field := range rules.RequiredFields {
		if err := validateField(parsedData, field); err != nil {
			return fmt.Errorf("validation error for %s: %w", field, err)
		}
	}

	return nil
}
