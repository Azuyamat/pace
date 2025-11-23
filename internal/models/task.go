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
	Alias           string
	Command         string
	Inputs          []string
	Outputs         []string
	DependsOn       []string
	Env             map[string]string
	Cache           bool
	WorkingDir      string
	Requires        []string
	Triggers        []string
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
	Args            *TaskArgs
	ExtraArgs       []string
	When            string
}
