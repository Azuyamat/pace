package command

import "github.com/azuyamat/pace/internal/config"

type Flag struct {
	Label string
	Value interface{}
}

func NewFlag(label string, value interface{}) *Flag {
	return &Flag{Label: label, Value: value}
}

type CommandContext struct {
	Cfg   *config.Config
	flags map[string]*Flag
}

func NewCommandContext(cfg *config.Config) *CommandContext {
	return &CommandContext{Cfg: cfg, flags: make(map[string]*Flag)}
}

func (ctx *CommandContext) SetFlag(name string, value *Flag) {
	ctx.flags[name] = NewFlag(name, value)
}

func (ctx *CommandContext) SetFlags(flags map[string]*Flag) {
	for name, value := range flags {
		ctx.flags[name] = value
	}
}

func (ctx *CommandContext) HasFlag(name string) bool {
	_, exists := ctx.flags[name]
	return exists
}

func (ctx *CommandContext) GetFlag(name string) (*Flag, bool) {
	value, exists := ctx.flags[name]
	return value, exists
}

func (ctx *CommandContext) GetFlagOr(name string, defaultValue interface{}) interface{} {
	if value, exists := ctx.flags[name]; exists {
		return value.Value
	}
	return defaultValue
}

func (ctx *CommandContext) GetStringFlag(name string) string {
	value, exists := ctx.flags[name]
	if !exists {
		return ""
	}
	return value.Value.(string)
}

func (ctx *CommandContext) GetIntFlag(name string) int {
	value, exists := ctx.flags[name]
	if !exists {
		return 0
	}
	return value.Value.(int)
}

func (ctx *CommandContext) GetBoolFlag(name string) bool {
	value, exists := ctx.flags[name]
	if !exists {
		return false
	}

	if boolVal, ok := value.Value.(bool); ok {
		return boolVal
	}

	if strVal, ok := value.Value.(string); ok {
		return strVal == "true" || strVal == "1"
	}

	return false
}

func (ctx *CommandContext) GetConfig() *config.Config {
	return ctx.Cfg
}
