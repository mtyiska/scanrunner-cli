package kubernetes

import (
	"fmt"
	"strings"

	"github.com/mtyiska/scanrunner/internal/fileparser"
	"github.com/mtyiska/scanrunner/internal/model"
)

// ValidateKubernetesManifest validates a parsed Kubernetes manifest for compliance and best practices.
func ValidateKubernetesManifest(parsedData map[string]interface{}, rules model.Rules) error {
	// Step 1: Validate Required Fields (from rules)
	for _, field := range rules.RequiredFields {
		if err := fileparser.ValidateField(parsedData, field); err != nil {
			return fmt.Errorf("missing or invalid required field: %s, error: %w", field, err)
		}
	}

	// Step 3: PodSecurity Checks
	if err := validatePodSecurity(parsedData); err != nil {
		return fmt.Errorf("PodSecurity validation failed: %w", err)
	}

	// Step 4: Network Policy Validation
	if err := validateNetworkPolicies(parsedData); err != nil {
		return fmt.Errorf("NetworkPolicy validation failed: %w", err)
	}

	return nil
}

// Helper function to retrieve a nested field from the manifest
func getField(data map[string]interface{}, path string) (interface{}, bool) {
	parts := strings.Split(path, ".")
	current := data

	for _, part := range parts {
		// Handle array notation (e.g., containers[])
		if strings.HasSuffix(part, "[]") {
			key := strings.TrimSuffix(part, "[]")
			if items, ok := current[key].([]interface{}); ok {
				// Return first item for simplicity; modify for multi-item validation
				if len(items) > 0 {
					current = items[0].(map[string]interface{})
					continue
				}
			}
			return nil, false
		}

		if value, ok := current[part].(map[string]interface{}); ok {
			current = value
		} else {
			return nil, false
		}
	}
	return current, true
}

// Helper function for PodSecurity validation
func validatePodSecurity(data map[string]interface{}) error {
	containersPath := "spec.containers"
	if containers, exists := getField(data, containersPath); exists {
		containerList, ok := containers.([]interface{})
		if !ok {
			return fmt.Errorf("containers field is not an array")
		}
		for _, container := range containerList {
			containerMap, ok := container.(map[string]interface{})
			if !ok {
				continue
			}
			if securityContext, exists := containerMap["securityContext"].(map[string]interface{}); exists {
				if runAsRoot, ok := securityContext["runAsNonRoot"].(bool); !ok || !runAsRoot {
					return fmt.Errorf("container must set securityContext.runAsNonRoot to true")
				}
			} else {
				return fmt.Errorf("missing securityContext in container spec")
			}
		}
	}
	return nil
}

// Helper function for Network Policy validation
func validateNetworkPolicies(data map[string]interface{}) error {
	// Check if the resource kind is a workload that might need a NetworkPolicy
	if kind, exists := data["kind"]; exists {
		kindStr, ok := kind.(string)
		if !ok {
			return fmt.Errorf("invalid kind field format")
		}

		// NetworkPolicy is not required for non-NetworkPolicy kinds
		if kindStr != "NetworkPolicy" {
			fmt.Println("Warning: No NetworkPolicy defined for the workload. Consider adding one for better security.")
			return nil
		}
	}

	// If kind is NetworkPolicy, validation passes
	return nil
}
