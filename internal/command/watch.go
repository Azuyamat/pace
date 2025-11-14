package command

import (
	"os"

	"azuyamat.dev/pace/internal/config"
	"azuyamat.dev/pace/internal/logger"
	"azuyamat.dev/pace/internal/runner"
)

func Watch(cfg *config.Config, args []string) {
	if len(args) < 1 {
		logger.Error("no task specified for watch mode")
		os.Exit(1)
	}

	taskName := args[0]

	if alias, exists := cfg.Aliases[taskName]; exists {
		taskName = alias
	}

	task, exists := cfg.Tasks[taskName]
	if !exists {
		logger.Error("task %q not found", taskName)
		os.Exit(1)
	}

	if len(task.Inputs) == 0 {
		logger.Error("task %q has no inputs defined for watching", taskName)
		os.Exit(1)
	}

	r := runner.NewRunner(cfg)
	w := runner.NewWatcher(r, taskName, task.Inputs)

	if err := w.Watch(); err != nil {
		logger.Error("Watch error: %v", err)
		os.Exit(1)
	}
}
