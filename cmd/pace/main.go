package main

import (
	"os"

	"azuyamat.dev/pace/internal/command"
	"azuyamat.dev/pace/internal/config"
	"azuyamat.dev/pace/internal/logger"
	"azuyamat.dev/pace/internal/runner"
)

func main() {
	cfg, err := config.ParseFile("config.pace")
	if err != nil {
		logger.Error("Error: %v", err)
		return
	}

	if len(os.Args) < 2 {
		if cfg.DefaultTask != "" {
			r := runner.NewRunner(cfg)
			if err := r.RunTask(cfg.DefaultTask); err != nil {
				logger.Error("Error running default task: %v", err)
				os.Exit(1)
			}
			return
		}
		command.Help()
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "run":
		command.Run(cfg, os.Args[2:])
	case "list":
		command.List(cfg, os.Args[2:])
	case "watch":
		command.Watch(cfg, os.Args[2:])
	case "help", "-h", "--help":
		command.Help()
	default:
		taskName := cmd
		if alias, exists := cfg.Aliases[cmd]; exists {
			taskName = alias
		}

		if _, exists := cfg.Tasks[taskName]; exists {
			r := runner.NewRunner(cfg)
			// Pass any extra arguments after the task name
			if err := r.RunTask(taskName, os.Args[2:]...); err != nil {
				logger.Error("Error running task: %v", err)
				os.Exit(1)
			}
		} else if cfg.DefaultTask != "" {
			r := runner.NewRunner(cfg)
			if err := r.RunTask(cfg.DefaultTask); err != nil {
				logger.Error("Error running default task: %v", err)
				os.Exit(1)
			}
		} else {
			logger.Error("Unknown command or task: %s", cmd)
			command.Help()
			os.Exit(1)
		}
	}
}
