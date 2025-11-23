package processing

import (
	"fmt"
	"strings"
	"time"

	"github.com/azuyamat/pace/internal/config/types"
)

type Validator struct {
	config *types.Config
	errors []error
}

func NewValidator(config *types.Config) *Validator {
	return &Validator{
		config: config,
		errors: make([]error, 0),
	}
}

func (v *Validator) Validate() error {
	v.validateTasks()
	v.validateHooks()
	v.validateHookReferences()
	v.validateTaskDependencies()
	v.validateConstants()
	v.validateAliases()
	v.validateTimeouts()
	v.validateRetry()

	if len(v.errors) > 0 {
		return v.combineErrors()
	}

	return nil
}

func (v *Validator) validateTasks() {
	for name, task := range v.config.Tasks {
		if task.Name == "" {
			v.addError(fmt.Errorf("task '%s' has no name", name))
		}

		if task.Command == "" {
			v.addError(fmt.Errorf("task '%s' has no command", name))
		}

		if task.Cache && len(task.Inputs) == 0 {
			v.addError(fmt.Errorf("task '%s' has cache enabled but no inputs specified", name))
		}
	}
}

func (v *Validator) validateHooks() {
	for name, hook := range v.config.Hooks {
		if hook.Name == "" {
			v.addError(fmt.Errorf("hook '%s' has no name", name))
		}

		if hook.Command == "" {
			v.addError(fmt.Errorf("hook '%s' has no command", name))
		}
	}
}

func (v *Validator) validateHookReferences() {
	for taskName, task := range v.config.Tasks {
		for _, hookName := range task.Requires {
			if _, exists := v.config.Hooks[hookName]; !exists {
				v.addError(fmt.Errorf("task '%s' references non-existent requires hook '%s'", taskName, hookName))
			}
		}

		for _, hookName := range task.Triggers {
			if _, exists := v.config.Hooks[hookName]; !exists {
				v.addError(fmt.Errorf("task '%s' references non-existent triggers hook '%s'", taskName, hookName))
			}
		}

		for _, hookName := range task.OnSuccess {
			if _, exists := v.config.Hooks[hookName]; !exists {
				v.addError(fmt.Errorf("task '%s' references non-existent on_success hook '%s'", taskName, hookName))
			}
		}

		for _, hookName := range task.OnFailure {
			if _, exists := v.config.Hooks[hookName]; !exists {
				v.addError(fmt.Errorf("task '%s' references non-existent on_failure hook '%s'", taskName, hookName))
			}
		}
	}
}

func (v *Validator) validateTaskDependencies() {
	for taskName, task := range v.config.Tasks {
		for _, depName := range task.DependsOn {
			if _, exists := v.config.Tasks[depName]; !exists {
				v.addError(fmt.Errorf("task '%s' depends on non-existent task '%s'", taskName, depName))
			}
		}
	}

	if v.hasCyclicDependencies() {
		v.addError(fmt.Errorf("cyclic task dependencies detected"))
	}
}

func (v *Validator) hasCyclicDependencies() bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for taskName := range v.config.Tasks {
		if v.isCyclic(taskName, visited, recStack) {
			return true
		}
	}

	return false
}

func (v *Validator) isCyclic(taskName string, visited, recStack map[string]bool) bool {
	visited[taskName] = true
	recStack[taskName] = true

	task, exists := v.config.Tasks[taskName]
	if !exists {
		return false
	}

	for _, dep := range task.DependsOn {
		if !visited[dep] {
			if v.isCyclic(dep, visited, recStack) {
				return true
			}
		} else if recStack[dep] {
			return true
		}
	}

	recStack[taskName] = false
	return false
}

func (v *Validator) validateConstants() {
	for name := range v.config.Constants {
		if name == "" {
			v.addError(fmt.Errorf("constant with empty name found"))
		}
	}

	for name := range v.config.Globals {
		if name == "" {
			v.addError(fmt.Errorf("global with empty name found"))
		}
	}
}

func (v *Validator) addError(err error) {
	v.errors = append(v.errors, err)
}

func (v *Validator) combineErrors() error {
	var messages []string
	for _, err := range v.errors {
		messages = append(messages, err.Error())
	}
	return fmt.Errorf("validation failed:\n  - %s", strings.Join(messages, "\n  - "))
}

func (v *Validator) validateAliases() {
	for alias, taskName := range v.config.Aliases {
		if _, exists := v.config.Tasks[taskName]; !exists {
			v.addError(fmt.Errorf("alias '%s' references non-existent task '%s'", alias, taskName))
		}
	}
}

func (v *Validator) validateTimeouts() {
	for name, task := range v.config.Tasks {
		if task.Timeout != "" {
			if _, err := time.ParseDuration(task.Timeout); err != nil {
				v.addError(fmt.Errorf("task '%s' has invalid timeout format '%s': %v", name, task.Timeout, err))
			}
		}
		if task.RetryDelay != "" {
			if _, err := time.ParseDuration(task.RetryDelay); err != nil {
				v.addError(fmt.Errorf("task '%s' has invalid retry_delay format '%s': %v", name, task.RetryDelay, err))
			}
		}
	}
}

func (v *Validator) validateRetry() {
	for name, task := range v.config.Tasks {
		if task.Retry < 0 {
			v.addError(fmt.Errorf("task '%s' has negative retry count: %d", name, task.Retry))
		}
	}
}
