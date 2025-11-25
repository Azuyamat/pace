package generator

import (
	"github.com/azuyamat/pace/internal/config"
	"github.com/azuyamat/pace/internal/models"
)

func NewUnknownGenerator() Generator {
	return newGenerator(generateUnknown, models.ProjectTypeUnknown)
}

func generateUnknown() (config.Config, error) {
	cfg := config.NewDefaultConfig()

	cfg.DefaultTask = "build"

	cfg.Tasks["build"] = models.Task{
		Name:        "build",
		Alias:       "b",
		Command:     "echo 'No build command specified'",
		Description: "Build the project",
	}

	cfg.Tasks["test"] = models.Task{
		Name:        "test",
		Alias:       "t",
		Command:     "echo 'No test command specified'",
		Description: "Run tests",
	}

	cfg.Tasks["lint"] = models.Task{
		Name:        "lint",
		Alias:       "l",
		Command:     "echo 'No lint command specified'",
		Description: "Lint code",
	}

	return *cfg, nil
}
