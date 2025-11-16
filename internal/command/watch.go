package command

import (
	"fmt"

	"azuyamat.dev/pace/internal/config"
	"azuyamat.dev/pace/internal/runner"
)

func init() {
	CommandRegistry.Register(watchCommand())
}

func watchCommand() *Command {
	return NewCommand("watch", "Watch a task's inputs and re-run it on changes").
		Arg(NewStringArg("task", "Name of the task to watch", true)).
		SetHandler(NewHandler(
			func(ctx *CommandContext, args *ValidatedArgs) error {
				taskName := args.String("task")
				return Watch(ctx.GetConfig(), []string{taskName})
			}))
}

func Watch(cfg *config.Config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no task name provided for watch command")
	}

	taskName := args[0]

	if alias, exists := cfg.Aliases[taskName]; exists {
		taskName = alias
	}

	task, exists := cfg.Tasks[taskName]
	if !exists {
		return fmt.Errorf("task %q not found in configuration", taskName)
	}

	if len(task.Inputs) == 0 {
		return fmt.Errorf("task %q has no inputs defined for watching", taskName)
	}

	r := runner.NewRunner(cfg)
	w := runner.NewWatcher(r, taskName, task.Inputs)

	if err := w.Watch(); err != nil {
		return err
	}

	return nil
}
