package main

import (
	"os"

	"github.com/azuyamat/pace/internal/command"
)

const ConfigFile = "config.pace"

func main() {
	command.RootCommand.Run(os.Args[1:])
}
