package command

import (
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
	if taskName == "" {
		taskName = config.DefaultTask
	}
	runner := runner.NewRunner(config)
	return runner.RunTask(taskName)
}
