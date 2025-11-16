package command

import "azuyamat.dev/pace/internal/logger"

func init() {
	CommandRegistry.Register(helpCommand())
}

func helpCommand() *Command {
	return NewCommand("help", "Display help information").
		Arg(NewStringArg("command", "Command to get help for", false)).
		SetHandler(NewHandler(
			func(ctx *CommandContext, args *ValidatedArgs) {
				commandName := args.StringOr("command", "")
				logger.Debug("Help command invoked with argument: %s", commandName)
				if commandName == "" {
					logger.Info("Available commands:")
					for _, cmd := range CommandRegistry.Commands() {
						logger.Info("  %s: %s", cmd.Label, cmd.Description)
					}
					logger.Info("Use 'help <command>' to get more information about a specific command.")
				} else {
					cmd, exists := CommandRegistry.GetCommand(commandName)
					if !exists {
						logger.Error("Command '%s' not found.", commandName)
						return
					}
					logger.Info("Help for command '%s':", cmd.Label)
					logger.Info("Description: %s", cmd.Description)
					if len(cmd.Args) > 0 {
						logger.Info("Arguments:")
						for _, arg := range cmd.Args {
							req := "optional"
							if arg.Required() {
								req = "required"
							}
							logger.Info("  %s (%s): %s", arg.Label(), req, arg.Description())
						}
					} else {
						logger.Info("This command has no arguments.")
					}
				}
			}))
}
