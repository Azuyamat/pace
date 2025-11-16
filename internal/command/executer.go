package command

import "azuyamat.dev/pace/internal/config"

func Execute(raw []string, cfg *config.Config) error {
	ctx := NewCommandContext(cfg)
	entryCommand := raw[0]
	args := raw[1:]
	command, _ := CommandRegistry.GetCommand(entryCommand)
	return command.Execute(ctx, args)
}
