package models

type Task struct {
	Name         string
	Command      string
	Inputs       []string
	Outputs      []string
	Dependencies []string
	Env          map[string]string
	Cache        bool
	WorkingDir   string
}
