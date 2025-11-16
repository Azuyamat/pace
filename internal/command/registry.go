package command

var CommandRegistry = NewRegistry()

type Registry struct {
	commands map[string]*Command
}

func NewRegistry() *Registry {
	return &Registry{commands: make(map[string]*Command)}
}

func (r *Registry) Register(cmd *Command) {
	r.commands[cmd.Label] = cmd
}

func (r *Registry) GetCommand(label string) (*Command, bool) {
	cmd, exists := r.commands[label]
	return cmd, exists
}

func (r *Registry) Commands() []*Command {
	cmds := make([]*Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}
