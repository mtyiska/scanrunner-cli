// rules_loader.go
package pkg

import (
	"fmt"
	"log"
	"os"

	"github.com/mtyiska/scanrunner/internal/model"
	"gopkg.in/yaml.v2"
)

// DefaultRules provides default values for custom-rules.yaml
func DefaultRules() model.Rules {
	return model.Rules{
		RequiredFields: model.AllowedPrefixes,
	}
}

// LoadRules loads and validates the custom-rules.yaml file
func LoadRules(path string) (model.Rules, error) {
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
		return model.Rules{}, fmt.Errorf("failed to parse rules file: %w", err)
	}

	// Validate the rules
	if err := model.ValidateRules(rules); err != nil {
		return model.Rules{}, fmt.Errorf("invalid rules file: %w", err)
	}

	// log.Printf("Rules loaded from %s: %+v\n", path, rules)
	return rules, nil
}

