// model/rules.go
package model

import (
	"fmt"
	"strings"
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



// validateRules checks the structure and syntax of the rules
func ValidateRules(rules Rules) error {

	seen := make(map[string]bool)
	for _, field := range rules.RequiredFields {
		// Check for duplicates
		if seen[field] {
			return fmt.Errorf("duplicate rule found: %s", field)
		}
		seen[field] = true

		// Check for valid syntax
		if err := ValidateFieldSyntax(field); err != nil {
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
func ValidateFieldSyntax(field string) error {
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
