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

func NewDefaultConfig() *Config {
	return &Config{
		Tasks:     make(map[string]models.Task),
		Hooks:     make(map[string]models.Hook),
		Globals:   make(map[string]string),
		Constants: make(map[string]string),
		Aliases:   make(map[string]string),
	}
}

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

	if err := processImports(cfg, filepath.Dir(path)); err != nil {
		return nil, err
	}

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

func (cfg *Config) GetTask(name string) (models.Task, bool) {
	task, exists := cfg.Tasks[name]
	return task, exists
}

func (cfg *Config) GetTaskOrDefault(name string) (models.Task, bool) {
	if name == "" && cfg.DefaultTask != "" {
		name = cfg.DefaultTask
	}
	return cfg.GetTask(name)
}

func (cfg *Config) GetHook(name string) (models.Hook, bool) {
	hook, exists := cfg.Hooks[name]
	return hook, exists
}

func processImports(cfg *Config, baseDir string) error {
	for _, importPath := range cfg.Imports {
		fullPath := filepath.Join(baseDir, importPath)
		importedCfg, err := ParseFile(fullPath)
		if err != nil {
			return err
		}

		importField(importedCfg.Tasks, cfg.Tasks)
		importField(importedCfg.Hooks, cfg.Hooks)
		importField(importedCfg.Constants, cfg.Constants)
		importField(importedCfg.Globals, cfg.Globals)
		importField(importedCfg.Aliases, cfg.Aliases)
	}

	return nil
}

func importField[T any](src, dest map[string]T) {
	for name, value := range src {
		if _, exists := dest[name]; !exists {
			dest[name] = value
		}
	}
}
