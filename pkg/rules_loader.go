// rules_loader.go
package pkg

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Rules represents the expected structure of custom-rules.yaml
type Rules struct {
	RequiredFields []string `yaml:"required_fields"` // List of required fields
}

// AllowedPrefixes defines the valid prefixes for required fields
var AllowedPrefixes = []string{
	// General fields
	"apiVersion",
	"kind",
	"metadata",
	"metadata.name",
	"spec",

	// Metadata fields
	"metadata.labels",
	"metadata.annotations",

	// Pod or workload fields
	"spec.containers",
	"spec.containers[].name",
	"spec.containers[].image",
	"spec.containers[].resources.limits",
	"spec.containers[].resources.requests",

	// Security fields
	"spec.securityContext",
	"spec.containers[].securityContext",
	"spec.serviceAccountName",

	// Deployment-specific fields
	"spec.replicas",
	"spec.selector",
	"spec.template.metadata",
	"spec.template.spec",
}

// DefaultRules provides default values for custom-rules.yaml
func DefaultRules() Rules {
	return Rules{
		RequiredFields: AllowedPrefixes,
	}
}

// LoadRules loads and validates the custom-rules.yaml file
func LoadRules(path string) (Rules, error) {
	// Default to the config/default-rules.yaml file if no path is provided
	if path == "" {
		path = "./config/default-rules.yaml"
	}

	rules := DefaultRules() // Start with defaults

	file, err := os.Open(path)
	if err != nil {
		log.Printf("Rules file not found at %s. Using defaults. Error: %v\n", path, err)
		return rules, nil // Return defaults if file is missing
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	decoder.SetStrict(false) // Ignore unknown fields

	if err := decoder.Decode(&rules); err != nil {
		return Rules{}, fmt.Errorf("failed to parse rules file: %w", err)
	}

	// Validate the rules
	if err := validateRules(rules); err != nil {
		return Rules{}, fmt.Errorf("invalid rules file: %w", err)
	}

	log.Printf("Rules loaded from %s: %+v\n", path, rules)
	return rules, nil
}

// validateRules checks the structure and syntax of the rules
func validateRules(rules Rules) error {

	seen := make(map[string]bool)
	for _, field := range rules.RequiredFields {
		// Check for duplicates
		if seen[field] {
			return fmt.Errorf("duplicate rule found: %s", field)
		}
		seen[field] = true

		// Check for valid syntax
		if err := validateFieldSyntax(field); err != nil {
			return fmt.Errorf("invalid rule syntax for field '%s': %w", field, err)
		}

		// Check for allowed prefixes
		valid := false
		for _, prefix := range AllowedPrefixes {
			if strings.HasPrefix(field, prefix) {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("unsupported field path: %s", field)
		}
	}
	return nil
}

// validateFieldSyntax ensures the field path uses valid syntax
func validateFieldSyntax(field string) error {
	if strings.TrimSpace(field) == "" {
		return fmt.Errorf("field path cannot be empty")
	}
	if strings.Contains(field, "..") {
		return fmt.Errorf("field path contains invalid sequence: '..'")
	}
	if strings.HasSuffix(field, ".") || strings.HasSuffix(field, "[") {
		return fmt.Errorf("field path has invalid ending: '%s'", field)
	}
	return nil
}
