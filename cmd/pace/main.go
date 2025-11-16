package main

import (
	"os"

	"azuyamat.dev/pace/internal/command"
	"azuyamat.dev/pace/internal/config"
	"azuyamat.dev/pace/internal/logger"
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
