package config

import (
	"os"

	"azuyamat.dev/pace/internal/models"
)

type Config struct {
	Tasks map[string]models.Task
}

func ParseFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Parse(string(data))
}
