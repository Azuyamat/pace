package main

import (
	"os"

	"azuyamat.dev/pace/internal/command"
	"azuyamat.dev/pace/internal/config"
	"azuyamat.dev/pace/internal/logger"
)

const CONFIG_FILE = "config.pace"

func main() {
	cfg, err := config.ParseFile(CONFIG_FILE)
	if err != nil {
		logger.Error("Error: %v", err)
		return
	}
	if err := command.Execute(os.Args[1:], cfg); err != nil {
		logger.Error("Error: %v", err)
	}
}
