package main

import (
	"os"

	"github.com/azuyamat/pace/internal/command"
	"github.com/azuyamat/pace/internal/config"
	"github.com/azuyamat/pace/internal/logger"
)

const ConfigFile = "config.pace"

func main() {
	cfg, err := config.ParseFile(ConfigFile)
	if err != nil {
		logger.Error("Error: %v", err)
		return
	}
	if err := command.Execute(os.Args[1:], cfg); err != nil {
		logger.Error("Error: %v", err)
	}
}
