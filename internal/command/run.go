package command

import (
	"fmt"

	gear "github.com/azuyamat/gear/command"
	"github.com/azuyamat/pace/internal/config"
	"github.com/azuyamat/pace/internal/runner"
)

var runCommand = gear.NewExecutableCommand("run", "Run a specified task").
	Args(
		gear.NewStringArg("task", "Name of the task to run").AsOptional()).
	Handler(runHandler)

func init() {
	RootCommand.AddChild(runCommand)
}

func runHandler(ctx *gear.Context, args gear.ValidatedArgs) error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}
	taskName := args.String("task")
	runner := runner.NewRunner(config)

	task, exists := config.GetTaskOrDefault(taskName)
	if !exists {
		return fmt.Errorf("task '%s' not found", taskName)
	}

	if task.Watch {
		return Watch(config, []string{taskName})
	}

	return runner.RunTask(task)
}
