package main

import (
	"os"

	"github.com/azuyamat/pace/internal/command"
	"github.com/azuyamat/pace/internal/logger"
)

func main() {
	err := command.RootCommand.Run(os.Args[1:])
	if err != nil {
		logger.Error("%s", err.Error())
		os.Exit(1)
	}
}
