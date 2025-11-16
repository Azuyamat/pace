package command

import (
	"fmt"

	"azuyamat.dev/pace/internal/config"
	"azuyamat.dev/pace/internal/logger"
)

func Execute(raw []string, cfg *config.Config) error {
	flags, args := extractFlags(raw)
	logger.Debug("Extracted flags: %v, args: %v", flags, args)
	ctx := NewCommandContext(cfg)
	if len(raw) == 0 {
		return executeCommand("help", []string{}, ctx)
	}
	entryCommand := args[0]
	for name, value := range flags {
		ctx.SetFlag(name, value)
	}
	args = args[1:]
	return executeCommand(entryCommand, args, ctx)
}

func executeCommand(cmdLabel string, rawArgs []string, ctx *CommandContext) error {
	command, exists := CommandRegistry.GetCommand(cmdLabel)
	if !exists {
		return fmt.Errorf("command '%s' not found", cmdLabel)
	}
	return command.Execute(ctx, rawArgs)
}

func extractFlags(rawArgs []string) (map[string]string, []string) {
	flags := make(map[string]string)
	args := []string{}
	skipNext := false

	for i, arg := range rawArgs {
		if skipNext {
			skipNext = false
			continue
		}
		if len(arg) > 2 && arg[0:2] == "--" {
			eqIndex := -1
			for j := 2; j < len(arg); j++ {
				if arg[j] == '=' {
					eqIndex = j
					break
				}
			}
			if eqIndex != -1 {
				flagName := arg[2:eqIndex]
				flagValue := arg[eqIndex+1:]
				flags[flagName] = flagValue
			} else if i+1 < len(rawArgs) {
				flagName := arg[2:]
				flagValue := rawArgs[i+1]
				flags[flagName] = flagValue
				skipNext = true
			}
		} else {
			args = append(args, arg)
		}
	}
	return flags, args
}
