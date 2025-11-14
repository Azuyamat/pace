package models

type Hook struct {
	Name        string
	Command     string
	Env         map[string]string
	WorkingDir  string
	Description string
}

type TaskArgs struct {
	Required []string
	Optional []string
}

type Task struct {
	Name            string
	Command         string
	Inputs          []string
	Outputs         []string
	Dependencies    []string
	Env             map[string]string
	Cache           bool
	WorkingDir      string
	BeforeHooks     []string
	AfterHooks      []string
	OnSuccess       []string
	OnFailure       []string
	Description     string
	Watch           bool
	Parallel        bool
	Silent          bool
	ContinueOnError bool
	Timeout         string
	Retry           int
	RetryDelay      string
	Args            *TaskArgs // Argument definition
	ExtraArgs       []string  // Additional arguments passed at runtime
}
