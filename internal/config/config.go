package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/azuyamat/pace/internal/config/loading"
)

type Config = loading.Config

var ConfigFile = loading.ConfigFile

func NewDefaultConfig() *Config {
	return loading.NewDefaultConfig()
}

func GetConfig() (*Config, error) {
	return loading.GetConfig()
}

func ParseFile(path string) (*Config, error) {
	return loading.ParseFile(path)
}

func UpdateGitignore(projectPath string) error {
	gitignorePath := fmt.Sprintf("%s/.gitignore", projectPath)

	content, err := os.ReadFile(gitignorePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	gitignoreContent := string(content)
	cacheEntry := ".pace-cache/"

	if strings.Contains(gitignoreContent, cacheEntry) {
		return nil
	}

	if len(gitignoreContent) > 0 && !strings.HasSuffix(gitignoreContent, "\n") {
		gitignoreContent += "\n"
	}

	gitignoreContent += cacheEntry + "\n"

	return os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
}
