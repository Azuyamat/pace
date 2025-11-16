package command

import "fmt"

type Handler func(ctx *CommandContext, args *ValidatedArgs) error

func NewHandler(f Handler) Handler {
	return f
}

type Command struct {
	Label       string
	Description string
	Subcommands []Command
	Args        Args
	Handler     Handler
}

func NewCommand(label, description string) *Command {
	return &Command{
		Label:       label,
		Description: description,
		Subcommands: []Command{},
		Args:        Args{},
		Handler: func(ctx *CommandContext, args *ValidatedArgs) error {
			return fmt.Errorf("no handler defined for command %q", label)
		},
	}
}

func (c *Command) Subcommand(cmd Command) *Command {
	c.Subcommands = append(c.Subcommands, cmd)
	return c
}

func (c *Command) Arg(arg Arg) *Command {
	c.Args = append(c.Args, arg)
	return c
}

func (c *Command) SetHandler(handler Handler) *Command {
	c.Handler = handler
	return c
}

func (c *Command) Execute(ctx *CommandContext, rawArgs []string) error {
	validated, err := c.validateAndParse(rawArgs)
	if err != nil {
		return err
	}

	if err = c.Handler(ctx, validated); err != nil {
		return err
	}

	return nil
}

func (c *Command) validateAndParse(rawArgs []string) (*ValidatedArgs, error) {
	validated := NewValidatedArgs()
	requiredCount := c.requiredArgsCount()

	if len(rawArgs) < requiredCount {
		return nil, fmt.Errorf("not enough arguments provided for command %q", c.Label)
	}

	for i, arg := range c.Args {
		if i >= len(rawArgs) {
			if arg.Required() {
				return nil, fmt.Errorf("missing required argument %q for command %q", arg.Label(), c.Label)
			}
			break
		}

		rawValue := rawArgs[i]

		switch a := arg.(type) {
		case *StringArg:
			validated.values[a.Label()] = rawValue
		case *IntArg:
			var intValue int
			_, err := fmt.Sscanf(rawValue, "%d", &intValue)
			if err != nil {
				return nil, fmt.Errorf("invalid integer value for argument %q: %v", a.Label(), err)
			}
			validated.values[a.Label()] = intValue
		case *BoolArg:
			var boolValue bool
			_, err := fmt.Sscanf(rawValue, "%t", &boolValue)
			if err != nil {
				return nil, fmt.Errorf("invalid boolean value for argument %q: %v", a.Label(), err)
			}
			validated.values[a.Label()] = boolValue
		default:
			return nil, fmt.Errorf("unknown argument type for %q", arg.Label())
		}
	}

	return validated, nil
}

func (c *Command) requiredArgsCount() int {
	count := 0
	for _, arg := range c.Args {
		if arg.Required() {
			count++
		}
	}

	return count
}
