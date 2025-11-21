package main

import (
	"fmt"
	"os"

	"github.com/azuyamat/pace/internal/command"
)

const ConfigFile = "config.pace"

func main() {
	err := command.RootCommand.Run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
