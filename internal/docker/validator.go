package docker

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser" // For parsing Dockerfiles
)

// ValidateDockerfile validates a Dockerfile for best practices, linting, and security checks.
func ValidateDockerfile(filePath string) error {
	// Step 1: Read the Dockerfile
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read Dockerfile: %w", err)
	}

	// Step 2: Parse and analyze the Dockerfile content
	parsedDockerfile, err := parseDockerfile(content)
	if err != nil {
		return fmt.Errorf("failed to parse Dockerfile: %w", err)
	}

	// Step 3: Perform linting checks
	if err := lintDockerfile(parsedDockerfile); err != nil {
		return fmt.Errorf("linting failed: %w", err)
	}

	// Step 4: Perform security scanning (Trivy)
	if err := scanDockerfileForSecrets(filePath); err != nil {
		return fmt.Errorf("security scan failed: %w", err)
	}

	// Step 5: Return success if all checks pass
	return nil
}

// parseDockerfile parses the Dockerfile content using the BuildKit parser.
func parseDockerfile(content []byte) (*parser.Node, error) {
	reader := strings.NewReader(string(content))
	parsed, err := parser.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Dockerfile: %w", err)
	}
	return parsed.AST, nil
}

// lintDockerfile performs linting and best practices validation on the parsed Dockerfile.
func lintDockerfile(ast *parser.Node) error {
	for _, child := range ast.Children {
		switch strings.ToUpper(child.Value) {
		case "ADD":
			return fmt.Errorf("line %d: use 'COPY' instead of 'ADD' for better security", child.StartLine)
		case "FROM":
			if len(child.Next.Value) == 0 || strings.Contains(child.Next.Value, "latest") {
				return fmt.Errorf("line %d: avoid using 'latest' tag in FROM directive for better reproducibility", child.StartLine)
			}
		case "RUN":
			if strings.Contains(child.Original, "apt-get install") && !strings.Contains(child.Original, "apt-get update") {
				return fmt.Errorf("line %d: missing 'apt-get update' before 'apt-get install'", child.StartLine)
			}
		}
	}
	return nil
}



// scanDockerfileForSecrets scans the Dockerfile for secrets using Trivy.
func scanDockerfileForSecrets(filePath string) error {
	// fmt.Printf("Scanning %s for secrets using Trivy...\n", filePath)

	// Build the Trivy command to scan the file
	cmd := exec.Command("trivy", "fs", "--security-checks", "secret", "--exit-code", "0", "--no-progress", filePath)

	// Capture the output and errors
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Run the command
    err := cmd.Run()
    if err != nil {
        if _, ok := err.(*exec.Error); ok {
            return fmt.Errorf("Trivy is not installed or not in PATH. Please install it and try again")
        }
        fmt.Printf("Error during Trivy scan: %s\n", stderr.String())
        return fmt.Errorf("Trivy scan failed for %s: %w", filePath, err)
    }
	// Print the scan results
	// fmt.Printf("Trivy scan results for %s:\n%s", filePath, out.String())
	return nil
}
