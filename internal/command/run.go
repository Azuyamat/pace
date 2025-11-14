package command

import (
	"flag"
	"os"

	"azuyamat.dev/pace/internal/config"
	"azuyamat.dev/pace/internal/logger"
	"azuyamat.dev/pace/internal/runner"
)

func Run(cfg *config.Config, args []string) {
	runFlags := flag.NewFlagSet("run", flag.ExitOnError)
	dryRun := runFlags.Bool("dry-run", false, "Show what would be executed without running")
	force := runFlags.Bool("force", false, "Force rebuild, ignoring cache")

	runFlags.Parse(args)

	if runFlags.NArg() < 1 {
		if cfg.DefaultTask != "" {
			r := runner.NewRunner(cfg)
			r.DryRun = *dryRun
			r.Force = *force
			if err := r.RunTask(cfg.DefaultTask); err != nil {
				logger.Error("Error running default task: %v", err)
				os.Exit(1)
			}
			return
		}
		logger.Error("no task specified and no default task set")
		os.Exit(1)
	}

	taskName := runFlags.Arg(0)

	if alias, exists := cfg.Aliases[taskName]; exists {
		taskName = alias
	}

	r := runner.NewRunner(cfg)
	r.DryRun = *dryRun
	r.Force = *force

	// Pass extra arguments (everything after the task name) to the task
	extraArgs := runFlags.Args()[1:]
	if err := r.RunTask(taskName, extraArgs...); err != nil {
		logger.Error("Error running task: %v", err)
		os.Exit(1)
	}
}
