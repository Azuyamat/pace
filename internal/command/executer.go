package command

import (
	"fmt"
	"strings"

	"azuyamat.dev/pace/internal/config"
	"azuyamat.dev/pace/internal/logger"
)

func Execute(raw []string, cfg *config.Config) error {
	flags, args := extractFlags(raw)
	logger.Debug("Extracted flags: %+v", flags)
	logger.Debug("Remaining args: %+v", args)
	ctx := NewCommandContext(cfg)
	if len(args) == 0 {
		return executeCommand("help", []string{}, ctx)
	}
	entryCommand := args[0]
	ctx.SetFlags(flags)
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

func extractFlags(rawArgs []string) (map[string]*Flag, []string) {
	flags := make(map[string]*Flag)
	args := []string{}

	for i := 0; i < len(rawArgs); i++ {
		arg := rawArgs[i]
		if !isFlag(arg) {
			args = append(args, arg)
			continue
		}

		flag, err := parseFlag(arg)
		if err != nil {
			args = append(args, arg)
			continue
		}

		flags[flag.Label] = flag
	}

	return flags, args
}

func isFlag(arg string) bool {
	return len(arg) > 0 && strings.HasPrefix(arg, "--")
}

func parseFlag(arg string) (*Flag, error) {
	if !isFlag(arg) {
		return nil, fmt.Errorf("argument %q is not a flag", arg)
	}

	var label, value string
	if eqIndex := strings.Index(arg, "="); eqIndex != -1 {
		label = arg[2:eqIndex]
		value = arg[eqIndex+1:]
	} else {
		label = arg[2:]
		value = "true"
	}

	return NewFlag(label, value), nil
}
