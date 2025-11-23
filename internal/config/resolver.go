package config

import (
	"os"
	"regexp"
	"strings"

	"github.com/azuyamat/pace/internal/logger"
)

var varPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

type Resolver struct {
	config         *Config
	unresolvedVars map[string]bool
}

func NewResolver(config *Config) *Resolver {
	return &Resolver{
		config:         config,
		unresolvedVars: make(map[string]bool),
	}
}

func (r *Resolver) ResolveString(input string) string {
	return varPattern.ReplaceAllStringFunc(input, func(match string) string {
		varName := match[2 : len(match)-1]

		if value, exists := r.config.Constants[varName]; exists {
			return value
		}

		if value, exists := r.config.Globals[varName]; exists {
			return value
		}

		if value := os.Getenv(varName); value != "" {
			return value
		}

		if !r.unresolvedVars[varName] {
			logger.Warning("Unresolved variable: ${%s}", varName)
			r.unresolvedVars[varName] = true
		}

		return match
	})
}

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
