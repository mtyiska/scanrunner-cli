package pkg

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


// validateField traverses the YAML structure to check if the required field exists
func validateField(data interface{}, fieldPath string) error {
	if strings.TrimSpace(fieldPath) == "" {
        return nil // Stop traversal if there are no more parts to validate
    }
    parts := strings.Split(fieldPath, ".")
    current := data

    for i, part := range parts {
        // Handle array notation (e.g., spec.containers[])
        if strings.HasSuffix(part, "[]") {
            partKey := strings.TrimSuffix(part, "[]")
            mapData, ok := current.(map[string]interface{})
            if !ok {
                fmt.Printf("ERROR: Expected map[string]interface{}, got %T\n", current)
                return fmt.Errorf("field '%s' not found in map", strings.Join(parts[:i], "."))
            }

            array, exists := mapData[partKey]
            if !exists {
                fmt.Printf("ERROR: Array key '%s' not found in map: %+v\n", partKey, mapData)
                return fmt.Errorf("missing required array field: %s", strings.Join(parts[:i+1], "."))
            }

            arrayItems, ok := array.([]interface{})
            if !ok {
                fmt.Printf("ERROR: Field '%s' is not an array. Value: %+v\n", partKey, array)
                return fmt.Errorf("field '%s' is not an array", strings.Join(parts[:i+1], "."))
            }

            for _, item := range arrayItems {
                if err := validateField(item, strings.Join(parts[i+1:], ".")); err != nil {
                    return fmt.Errorf("error in array '%s': %v", strings.Join(parts[:i+1], "."), err)
                }
            }
            return nil // Successfully validated all array items
        }

        // Handle regular fields
        mapData, ok := current.(map[string]interface{})
        if !ok {
            fmt.Printf("ERROR: Expected map[string]interface{}, got %T\n", current)
            return fmt.Errorf("field '%s' not found in map", strings.Join(parts[:i], "."))
        }

        value, exists := mapData[part]
        if !exists {
            fmt.Printf("ERROR: Key '%s' not found in current map: %+v\n", part, mapData)
            return fmt.Errorf("missing required field: %s", strings.Join(parts[:i+1], "."))
        }

        current = value // Traverse deeper
    }

    return nil // All parts of the fieldPath validated successfully
}
