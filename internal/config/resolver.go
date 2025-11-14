package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var varPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

type Resolver struct {
	config *Config
}

func NewResolver(config *Config) *Resolver {
	return &Resolver{config: config}
}

// ResolveString resolves variable references in a string
func (r *Resolver) ResolveString(input string) string {
	return varPattern.ReplaceAllStringFunc(input, func(match string) string {
		// Extract variable name from ${VAR}
		varName := match[2 : len(match)-1]

		// Check in constants first
		if value, exists := r.config.Constants[varName]; exists {
			return value
		}

		// Check in globals
		if value, exists := r.config.Globals[varName]; exists {
			return value
		}

		// Check environment variables
		if value := os.Getenv(varName); value != "" {
			return value
		}

		// Return original if not found
		return match
	})
}

// ResolveStringSlice resolves variables in a slice of strings
func (r *Resolver) ResolveStringSlice(slice []string) []string {
	result := make([]string, len(slice))
	for i, s := range slice {
		result[i] = r.ResolveString(s)
	}
	return result
}

// ResolveStringMap resolves variables in map values
func (r *Resolver) ResolveStringMap(m map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = r.ResolveString(v)
	}
	return result
}

// ExpandEnvVars expands environment variables in format $VAR or ${VAR}
func ExpandEnvVars(s string) string {
	// Handle ${VAR} format
	s = varPattern.ReplaceAllStringFunc(s, func(match string) string {
		varName := match[2 : len(match)-1]
		if value := os.Getenv(varName); value != "" {
			return value
		}
		return match
	})

	// Handle $VAR format (simple case)
	parts := strings.Split(s, "$")
	if len(parts) <= 1 {
		return s
	}

	var result strings.Builder
	result.WriteString(parts[0])

	for i := 1; i < len(parts); i++ {
		part := parts[i]
		if len(part) == 0 {
			result.WriteString("$")
			continue
		}

		// Extract variable name (alphanumeric and underscore)
		varEnd := 0
		for varEnd < len(part) && (isLetter(part[varEnd]) || isDigit(part[varEnd])) {
			varEnd++
		}

		if varEnd > 0 {
			varName := part[:varEnd]
			if value := os.Getenv(varName); value != "" {
				result.WriteString(value)
				result.WriteString(part[varEnd:])
			} else {
				result.WriteString("$")
				result.WriteString(part)
			}
		} else {
			result.WriteString("$")
			result.WriteString(part)
		}
	}

	return result.String()
}

// EvaluateCondition evaluates simple conditional expressions
func EvaluateCondition(condition string) (bool, error) {
	condition = strings.TrimSpace(condition)

	// Handle OS checks: OS == "windows" or OS == "linux"
	if strings.Contains(condition, "OS") {
		osPattern := regexp.MustCompile(`OS\s*==\s*"([^"]+)"`)
		matches := osPattern.FindStringSubmatch(condition)
		if len(matches) > 1 {
			targetOS := strings.ToLower(matches[1])
			currentOS := strings.ToLower(os.Getenv("GOOS"))
			if currentOS == "" {
				currentOS = "windows" // Default for this system
			}
			return currentOS == targetOS, nil
		}
	}

	// Handle environment variable checks: ENV_VAR == "value"
	envPattern := regexp.MustCompile(`([A-Z_][A-Z0-9_]*)\s*==\s*"([^"]+)"`)
	matches := envPattern.FindStringSubmatch(condition)
	if len(matches) > 2 {
		envVar := matches[1]
		expectedValue := matches[2]
		actualValue := os.Getenv(envVar)
		return actualValue == expectedValue, nil
	}

	// Handle boolean environment variables: ENV_VAR
	if matched, _ := regexp.MatchString(`^[A-Z_][A-Z0-9_]*$`, condition); matched {
		value := os.Getenv(condition)
		return value != "" && value != "0" && strings.ToLower(value) != "false", nil
	}

	return false, fmt.Errorf("unable to evaluate condition: %s", condition)
}
