package template

import (
	"github.com/azuyamat/pace/internal/config"
	"github.com/azuyamat/pace/internal/models"
)

func NewGoGenerator() Generator {
	return newGenerator(generate, models.ProjectTypeGo)
}

func generate() (config.Config, error) {
	cfg := config.NewDefaultConfig()

	cfg.Tasks["build"] = models.Task{
		Name:    "build",
		Command: "go build -o bin/main ./...",
		Inputs:  []string{"**/*.go"},
		Outputs: []string{"bin/main"},
	}
	return *cfg, nil
}
