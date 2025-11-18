package config

import (
	"os"
	"path/filepath"

	"github.com/azuyamat/pace/internal/models"
)

type Config struct {
	Tasks       map[string]models.Task
	Hooks       map[string]models.Hook
	Globals     map[string]string
	Constants   map[string]string
	DefaultTask string
	Aliases     map[string]string
	Imports     []string
}

var ConfigFile = "config.pace"

func GetConfig() (*Config, error) {
	return ParseFile(ConfigFile)
}

func ParseFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg, err := Parse(string(data))
	if err != nil {
		return nil, err
	}

	// Process imports
	if len(cfg.Imports) > 0 {
		baseDir := filepath.Dir(path)
		for _, importPath := range cfg.Imports {
			fullPath := filepath.Join(baseDir, importPath)
			importedCfg, err := ParseFile(fullPath)
			if err != nil {
				return nil, err
			}

			// Merge imported config
			for name, task := range importedCfg.Tasks {
				if _, exists := cfg.Tasks[name]; !exists {
					cfg.Tasks[name] = task
				}
			}
			for name, hook := range importedCfg.Hooks {
				if _, exists := cfg.Hooks[name]; !exists {
					cfg.Hooks[name] = hook
				}
			}
			for name, value := range importedCfg.Constants {
				if _, exists := cfg.Constants[name]; !exists {
					cfg.Constants[name] = value
				}
			}
			for name, value := range importedCfg.Globals {
				if _, exists := cfg.Globals[name]; !exists {
					cfg.Globals[name] = value
				}
			}
			for name, value := range importedCfg.Aliases {
				if _, exists := cfg.Aliases[name]; !exists {
					cfg.Aliases[name] = value
				}
			}
		}
	}

	// Resolve variables
	resolver := NewResolver(cfg)
	for name, task := range cfg.Tasks {
		task.Command = resolver.ResolveString(task.Command)
		task.Inputs = resolver.ResolveStringSlice(task.Inputs)
		task.Outputs = resolver.ResolveStringSlice(task.Outputs)
		task.WorkingDir = resolver.ResolveString(task.WorkingDir)
		task.Env = resolver.ResolveStringMap(task.Env)
		cfg.Tasks[name] = task
	}
	for name, hook := range cfg.Hooks {
		hook.Command = resolver.ResolveString(hook.Command)
		hook.WorkingDir = resolver.ResolveString(hook.WorkingDir)
		hook.Env = resolver.ResolveStringMap(hook.Env)
		cfg.Hooks[name] = hook
	}

	validator := NewValidator(cfg)
	if err := validator.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
