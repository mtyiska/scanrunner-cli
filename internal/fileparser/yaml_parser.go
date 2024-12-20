// yaml_parser.go
package fileparser

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// ConvertMapKeys recursively converts map[interface{}]interface{} to map[string]interface{}
func ConvertMapKeys(data interface{}) interface{} {
	switch v := data.(type) {
	case map[interface{}]interface{}:
		newMap := make(map[string]interface{})
		for key, value := range v {
			strKey := fmt.Sprintf("%v", key) // Convert key to string
			newMap[strKey] = ConvertMapKeys(value)
		}
		return newMap
	case []interface{}:
		for i, item := range v {
			v[i] = ConvertMapKeys(item)
		}
	}
	return data
}


// ParseAndConvertYAML loads a YAML file and converts it to map[string]interface{}
func ParseAndConvertYAML(filePath string) (map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open YAML file: %w", err)
	}
	defer file.Close()

	var parsedData map[interface{}]interface{}
	decoder := yaml.NewDecoder(file)
	decoder.SetStrict(false) // Allow extra fields
	if err := decoder.Decode(&parsedData); err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}

	// Convert map[interface{}]interface{} to map[string]interface{}
	convertedData := ConvertMapKeys(parsedData).(map[string]interface{})
	return convertedData, nil
}

// GetField retrieves a nested field from a parsed YAML structure based on a dot-separated path.
func GetField(data map[string]interface{}, path string) (interface{}, error) {
	parts := strings.Split(path, ".")
	current := data

	for _, part := range parts {
		// Handle array notation (e.g., containers[])
		if strings.HasSuffix(part, "[]") {
			key := strings.TrimSuffix(part, "[]")
			array, exists := current[key]
			if !exists {
				return nil, fmt.Errorf("array field '%s' not found", key)
			}

			// Check if the value is an array
			arrayItems, ok := array.([]interface{})
			if !ok {
				return nil, fmt.Errorf("field '%s' is not an array", key)
			}

			if len(arrayItems) > 0 {
				// Use the first item in the array for simplicity
				if itemMap, ok := arrayItems[0].(map[string]interface{}); ok {
					current = itemMap
					continue
				}
				return nil, fmt.Errorf("array items in '%s' are not valid objects", key)
			}
			return nil, fmt.Errorf("array '%s' is empty", key)
		}

		// Traverse to the next level
		value, exists := current[part]
		if !exists {
			return nil, fmt.Errorf("field '%s' not found", part)
		}

		// Update current level if it's a map
		if mapValue, ok := value.(map[string]interface{}); ok {
			current = mapValue
		} else {
			return value, nil // Return the value if it's not a map
		}
	}

	return current, nil
}

// ValidateField traverses the YAML structure to check if the required field exists
func ValidateField(data map[string]interface{}, fieldPath string) error {
	if _, err := GetField(data, fieldPath); err != nil {
		return err
	}
	return nil
}
