package loading

import (
	"os"
	"path/filepath"

	"github.com/azuyamat/pace/internal/config/parsing"
	"github.com/azuyamat/pace/internal/config/processing"
	"github.com/azuyamat/pace/internal/config/types"
)

type Config = types.Config

var ConfigFile = "config.pace"

func NewDefaultConfig() *Config {
	return types.NewConfig()
}

func GetConfig() (*Config, error) {
	return ParseFile(ConfigFile)
}

func ParseFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg, err := parsing.Parse(string(data))
	if err != nil {
		return nil, err
	}

	if err := processImports(cfg, filepath.Dir(path)); err != nil {
		return nil, err
	}

	resolver := processing.NewResolver(cfg)
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

	validator := processing.NewValidator(cfg)
	if err := validator.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
