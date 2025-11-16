package command

import "azuyamat.dev/pace/internal/config"

func Execute(raw []string, cfg *config.Config) error {
	if len(raw) == 0 {
		return executeCommand("help", []string{}, cfg)
	}
	entryCommand := raw[0]
	args := raw[1:]
	return executeCommand(entryCommand, args, cfg)
}

func executeCommand(cmdLabel string, rawArgs []string, cfg *config.Config) error {
	ctx := NewCommandContext(cfg)
	command, exists := CommandRegistry.GetCommand(cmdLabel)
	if !exists {
		return nil
	}
	return command.Execute(ctx, rawArgs)
}
