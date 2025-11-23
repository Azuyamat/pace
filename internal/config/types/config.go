package types

import "github.com/azuyamat/pace/internal/models"

type Config struct {
	Tasks       map[string]models.Task
	Hooks       map[string]models.Hook
	Globals     map[string]string
	Constants   map[string]string
	DefaultTask string
	Aliases     map[string]string
	Imports     []string
}

func NewConfig() *Config {
	return &Config{
		Tasks:     make(map[string]models.Task),
		Hooks:     make(map[string]models.Hook),
		Globals:   make(map[string]string),
		Constants: make(map[string]string),
		Aliases:   make(map[string]string),
		Imports:   make([]string, 0),
	}
}

func (cfg *Config) GetTask(name string) (models.Task, bool) {
	task, exists := cfg.Tasks[name]
	return task, exists
}

func (cfg *Config) GetTaskOrDefault(name string) (models.Task, bool) {
	if name == "" && cfg.DefaultTask != "" {
		name = cfg.DefaultTask
	}
	if alias, exists := cfg.Aliases[name]; exists {
		name = alias
	}
	return cfg.GetTask(name)
}

func (cfg *Config) GetHook(name string) (models.Hook, bool) {
	hook, exists := cfg.Hooks[name]
	return hook, exists
}
