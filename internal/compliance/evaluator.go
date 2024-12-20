package compliance

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mtyiska/scanrunner/internal/docker"
	"github.com/mtyiska/scanrunner/internal/fileparser"
	"github.com/mtyiska/scanrunner/internal/kubernetes"
	"github.com/mtyiska/scanrunner/internal/model"
)


func ValidateFile(filePath string, rules model.Rules) error {
	ext := strings.ToLower(filepath.Ext(filePath))
	fileName := strings.ToLower(filepath.Base(filePath))

	switch {
	case ext == ".yaml" || ext == ".yml":
		parsedData, err := fileparser.ParseAndConvertYAML(filePath)
		if err != nil {
			return fmt.Errorf("error parsing YAML file: %w", err)
		}
		if err := kubernetes.ValidateKubernetesManifest(parsedData, rules); err != nil {
			return fmt.Errorf("Kubernetes manifest validation failed: %w", err)
		}
		return nil

	case strings.Contains(fileName, "docker"): // Handle Docker-related files
		if err := docker.ValidateDockerfile(filePath); err != nil {
			return fmt.Errorf("Dockerfile validation failed: %w", err)
		}
		return nil

	default:
		return fmt.Errorf("unsupported file type: %s. Supported types are: .yaml, .yml, and Docker-related files", filePath)
	}
}