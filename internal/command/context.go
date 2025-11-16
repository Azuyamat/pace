package command

import "azuyamat.dev/pace/internal/config"

type CommandContext struct {
	Cfg   *config.Config
	flags map[string]interface{}
}

func NewCommandContext(cfg *config.Config) *CommandContext {
	return &CommandContext{Cfg: cfg, flags: make(map[string]interface{})}
}

func (ctx *CommandContext) SetFlag(name string, value interface{}) {
	ctx.flags[name] = value
}

func (ctx *CommandContext) GetFlag(name string) (interface{}, bool) {
	value, exists := ctx.flags[name]
	return value, exists
}

func (ctx *CommandContext) GetFlagOr(name string, defaultValue interface{}) interface{} {
	if value, exists := ctx.flags[name]; exists {
		return value
	}
	return defaultValue
}

func (ctx *CommandContext) GetConfig() *config.Config {
	return ctx.Cfg
}
