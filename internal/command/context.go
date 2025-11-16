package command

import "azuyamat.dev/pace/internal/config"

type CommandContext struct {
	Cfg *config.Config
}

func NewCommandContext(cfg *config.Config) *CommandContext {
	return &CommandContext{Cfg: cfg}
}
