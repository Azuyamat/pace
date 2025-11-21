package command

import (
	"fmt"

	gear "github.com/azuyamat/gear/command"
	"github.com/azuyamat/pace/internal/config"
	"github.com/azuyamat/pace/internal/runner"
)

var watchCommand = gear.NewExecutableCommand("watch", "Watch a task's inputs and re-run it on changes").
	Args(
		gear.NewStringArg("task", "Name of the task to watch")).
	Handler(watchHandler)

func init() {
	RootCommand.AddChild(watchCommand)
}

func watchHandler(ctx *gear.Context, args gear.ValidatedArgs) error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}
	taskName := args.String("task")
	return Watch(config, []string{taskName})
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
	w := runner.NewWatcher(r, task, task.Inputs)

	if err := w.Watch(); err != nil {
		return err
	}

	return nil
}
