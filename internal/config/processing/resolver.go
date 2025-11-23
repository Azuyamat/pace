package processing

import (
	"os"
	"regexp"

	"github.com/azuyamat/pace/internal/config/types"
	"github.com/azuyamat/pace/internal/logger"
)

var varPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

type Resolver struct {
	config         *types.Config
	unresolvedVars map[string]bool
}

func NewResolver(config *types.Config) *Resolver {
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

func (r *Resolver) ResolveStringMap(m map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = r.ResolveString(v)
	}
	return result
}

func ExpandEnvVars(s string) string {
	return varPattern.ReplaceAllStringFunc(s, func(match string) string {
		varName := match[2 : len(match)-1]
		if value := os.Getenv(varName); value != "" {
			return value
		}
		return match
	})
}
