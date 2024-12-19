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

// // validateField traverses the YAML structure to check if the required field exists
// func validateField(data interface{}, fieldPath string) error {
// 	if strings.TrimSpace(fieldPath) == "" {
//         return nil // Stop traversal if there are no more parts to validate
//     }
//     parts := strings.Split(fieldPath, ".")
//     current := data

//     for i, part := range parts {
//         // Handle array notation (e.g., spec.containers[])
//         if strings.HasSuffix(part, "[]") {
//             partKey := strings.TrimSuffix(part, "[]")
//             mapData, ok := current.(map[string]interface{})
//             if !ok {
//                 fmt.Printf("ERROR: Expected map[string]interface{}, got %T\n", current)
//                 return fmt.Errorf("field '%s' not found in map", strings.Join(parts[:i], "."))
//             }

//             array, exists := mapData[partKey]
//             if !exists {
//                 fmt.Printf("ERROR: Array key '%s' not found in map: %+v\n", partKey, mapData)
//                 return fmt.Errorf("missing required array field: %s", strings.Join(parts[:i+1], "."))
//             }

//             arrayItems, ok := array.([]interface{})
//             if !ok {
//                 fmt.Printf("ERROR: Field '%s' is not an array. Value: %+v\n", partKey, array)
//                 return fmt.Errorf("field '%s' is not an array", strings.Join(parts[:i+1], "."))
//             }

//             for _, item := range arrayItems {
//                 if err := validateField(item, strings.Join(parts[i+1:], ".")); err != nil {
//                     return fmt.Errorf("error in array '%s': %v", strings.Join(parts[:i+1], "."), err)
//                 }
//             }
//             return nil // Successfully validated all array items
//         }

//         // Handle regular fields
//         mapData, ok := current.(map[string]interface{})
//         if !ok {
//             fmt.Printf("ERROR: Expected map[string]interface{}, got %T\n", current)
//             return fmt.Errorf("field '%s' not found in map", strings.Join(parts[:i], "."))
//         }

//         value, exists := mapData[part]
//         if !exists {
//             fmt.Printf("ERROR: Key '%s' not found in current map: %+v\n", part, mapData)
//             return fmt.Errorf("missing required field: %s", strings.Join(parts[:i+1], "."))
//         }

//         current = value // Traverse deeper
//     }

//     return nil // All parts of the fieldPath validated successfully
// }




