package schema

// SchemaVersion is the compatibility version for the Pace language schema.
const SchemaVersion = 1

// LanguageSchema describes Pace syntax metadata consumed by parser-adjacent tooling,
// editor integrations, and docs generators.
type LanguageSchema struct {
	Version   int         `json:"version"`
	TopLevel  []Statement `json:"topLevel"`
	TaskProps []Property  `json:"taskProperties"`
	HookProps []Property  `json:"hookProperties"`
	ArgsProps []Property  `json:"argsProperties"`
}

type Statement struct {
	Name        string `json:"name"`
	Snippet     string `json:"snippet"`
	Description string `json:"description"`
	Deprecated  bool   `json:"deprecated,omitempty"`
}

type Property struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Snippet     string `json:"snippet"`
	Description string `json:"description"`
	Field       string `json:"field,omitempty"`
	AliasOf     string `json:"aliasOf,omitempty"`
	Deprecated  bool   `json:"deprecated,omitempty"`
}

// Get returns the Go-owned Pace language schema.
func Get() LanguageSchema {
	return LanguageSchema{
		Version: SchemaVersion,
		TopLevel: []Statement{
			{Name: "var", Snippet: "var ${1:VAR_NAME} = \"${2:value}\"", Description: "Defines a variable for interpolation in strings. Syntax: `var NAME = \"value\"`."},
			{Name: "default", Snippet: "default ${1:task_name}", Description: "Sets the default task that runs when Pace is invoked without a task name. Syntax: `default task_name`."},
			{Name: "alias", Snippet: "alias ${1:short} ${2:task_name}", Description: "Defines a standalone task alias. Supported, but inline aliases (`task build [b] { ... }`) are preferred.", Deprecated: true},
			{Name: "import", Snippet: "import \"${1:path/to/file.pace}\"", Description: "Imports another Pace config file. Syntax: `import \"path/to/file.pace\"`."},
			{Name: "hook", Snippet: "hook ${1:hook_name} {\n\tdescription \"${2:description}\"\n\tcommand \"${3:command}\"\n}", Description: "Defines a reusable hook that can run before, after, or on task success/failure. Syntax: `hook name { ... }`."},
			{Name: "task", Snippet: "task ${1:task_name} [${2:alias}] {\n\tdescription \"${3:description}\"\n\tcommand \"${4:command}\"\n}", Description: "Defines a runnable Pace task. Syntax: `task name [alias] { ... }`."},
		},
		TaskProps: []Property{
			{Name: "description", Type: "string", Field: "Description", Snippet: "description \"${1:task description}\"", Description: "Task description."},
			{Name: "command", Type: "string", Field: "Command", Snippet: "command \"${1:command to run}\"", Description: "Command to execute."},
			{Name: "depends-on", Type: "stringArray", Field: "DependsOn", AliasOf: "dependencies", Snippet: "depends-on [${1:build, test}]", Description: "Task dependencies (alias of dependencies)."},
			{Name: "dependencies", Type: "stringArray", Field: "DependsOn", AliasOf: "depends-on", Snippet: "dependencies [${1:build, test}]", Description: "Task dependencies (alias of depends-on)."},
			{Name: "requires", Type: "stringArray", Field: "Requires", AliasOf: "before", Snippet: "requires [${1:setup}]", Description: "Hooks to run before task (preferred; alias of before)."},
			{Name: "before", Type: "stringArray", Field: "Requires", AliasOf: "requires", Snippet: "before [${1:setup}]", Description: "Hooks to run before task (alias of requires)."},
			{Name: "triggers", Type: "stringArray", Field: "Triggers", AliasOf: "after", Snippet: "triggers [${1:cleanup}]", Description: "Hooks to run after task (preferred; alias of after)."},
			{Name: "after", Type: "stringArray", Field: "Triggers", AliasOf: "triggers", Snippet: "after [${1:cleanup}]", Description: "Hooks to run after task (alias of triggers)."},
			{Name: "on_success", Type: "stringArray", Field: "OnSuccess", Snippet: "on_success [${1:notify}]", Description: "Hooks to run on success."},
			{Name: "on_failure", Type: "stringArray", Field: "OnFailure", Snippet: "on_failure [${1:notify}]", Description: "Hooks to run on failure."},
			{Name: "inputs", Type: "stringArray", Field: "Inputs", Snippet: "inputs [${1:\"src/**/*.go\"}]", Description: "Input file patterns."},
			{Name: "outputs", Type: "stringArray", Field: "Outputs", Snippet: "outputs [${1:\"build/output\"}]", Description: "Output file patterns."},
			{Name: "env", Type: "stringMap", Field: "Env", Snippet: "env {\n\t${1:KEY} = \"${2:value}\"\n}", Description: "Environment variables."},
			{Name: "args", Type: "argsBlock", Snippet: "args {\n\trequired [${1:arg1}]\n}", Description: "Command arguments."},
			{Name: "cache", Type: "boolean", Field: "Cache", Snippet: "cache ${1|true,false|}", Description: "Enable caching (true/false)."},
			{Name: "parallel", Type: "boolean", Field: "Parallel", Snippet: "parallel ${1|true,false|}", Description: "Run dependencies in parallel (true/false)."},
			{Name: "silent", Type: "boolean", Field: "Silent", Snippet: "silent ${1|true,false|}", Description: "Suppress output (true/false)."},
			{Name: "watch", Type: "boolean", Field: "Watch", Snippet: "watch ${1|true,false|}", Description: "Enable watch mode (true/false)."},
			{Name: "continue_on_error", Type: "boolean", Field: "ContinueOnError", Snippet: "continue_on_error ${1|true,false|}", Description: "Continue on error (true/false)."},
			{Name: "timeout", Type: "string", Field: "Timeout", Snippet: "timeout \"${1:5m}\"", Description: "Execution timeout."},
			{Name: "retry", Type: "number", Field: "Retry", Snippet: "retry ${1:2}", Description: "Number of retries."},
			{Name: "retry_delay", Type: "string", Field: "RetryDelay", Snippet: "retry_delay \"${1:3s}\"", Description: "Delay between retries."},
			{Name: "working_dir", Type: "string", Field: "WorkingDir", Snippet: "working_dir \"${1:.}\"", Description: "Working directory for the task command."},
			{Name: "when", Type: "string", Field: "When", Snippet: "when \"${1:condition}\"", Description: "Condition that controls whether the task runs."},
		},
		HookProps: []Property{
			{Name: "description", Type: "string", Field: "Description", Snippet: "description \"${1:hook description}\"", Description: "Hook description."},
			{Name: "command", Type: "string", Field: "Command", Snippet: "command \"${1:command to run}\"", Description: "Command to execute."},
			{Name: "env", Type: "stringMap", Field: "Env", Snippet: "env {\n\t${1:KEY} = \"${2:value}\"\n}", Description: "Environment variables."},
			{Name: "working_dir", Type: "string", Field: "WorkingDir", Snippet: "working_dir \"${1:.}\"", Description: "Working directory for the hook command."},
		},
		ArgsProps: []Property{
			{Name: "required", Type: "stringArray", Snippet: "required [${1:arg1}]", Description: "Required arguments."},
			{Name: "optional", Type: "stringArray", Snippet: "optional [${1:arg1}]", Description: "Optional arguments."},
		},
	}
}
