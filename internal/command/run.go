package command

import (
	"azuyamat.dev/pace/internal/runner"
)

func init() {
	CommandRegistry.Register(runCommand())
}

func runCommand() *Command {
	return NewCommand("run", "Run a specified task").
		Arg(NewStringArg("task", "Name of the task to run", true)).
		SetHandler(NewHandler(
			func(ctx *CommandContext, args *ValidatedArgs) error {
				taskName := args.String("task")

				runner := runner.NewRunner(ctx.GetConfig())
				runner.DryRun = ctx.GetBoolFlag("dry-run")
				runner.Force = ctx.GetBoolFlag("force")

				if err := runner.RunTask(taskName); err != nil {
					return err
				}
				return nil
			}))
}
