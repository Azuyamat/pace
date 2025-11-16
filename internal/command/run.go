package command

import (
	"os"

	"azuyamat.dev/pace/internal/logger"
	"azuyamat.dev/pace/internal/runner"
)

func init() {
	CommandRegistry.Register(runCommand())
}

func runCommand() *Command {
	return NewCommand("run", "Run a specified task").
		Arg(NewStringArg("task", "Name of the task to run", true)).
		SetHandler(NewHandler(
			func(ctx *CommandContext, args *ValidatedArgs) {
				taskName := args.String("task")

				runner := runner.NewRunner(ctx.GetConfig())
				runner.DryRun = ctx.GetFlagOr("dry-run", false).(bool)
				runner.Force = ctx.GetFlagOr("force", false).(bool)

				if err := runner.RunTask(taskName); err != nil {
					logger.Error("Error running task: %v", err)
					os.Exit(1)
				}
			}))
}
