package main

import (
	"fmt"
	"os"

	"azuyamat.dev/pace/internal/config"
	"azuyamat.dev/pace/internal/runner"
)

func main() {
	// Parse configuration file
	config, err := config.ParseFile("config.pace")
	if err != nil {
		fmt.Println("Error parsing config file:", err)
		return
	}

	args := os.Args[1:]
	operation := args[0]

	if operation == "run" {
		runner := runner.NewRunner(config)
		taskName := args[1]
		err := runner.RunTask(taskName)
		if err != nil {
			fmt.Println("Error running task:", err)
			return
		}
	}
}
