package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"azuyamat.dev/pace/internal/models"
)

type Executor struct {
	shell  *Shell
	log    taskLogger
	DryRun bool
}

type taskLogger interface {
	Task(format string, args ...interface{})
	Success(format string, args ...interface{})
	Warning(format string, args ...interface{})
}

func NewExecutor(shell *Shell, log taskLogger, dryRun bool) *Executor {
	return &Executor{
		shell:  shell,
		log:    log,
		DryRun: dryRun,
	}
}

// interpolateArgs replaces argument placeholders in the command string
// Supports:
//
//	$@ - all arguments
//	$1, $2, $3, etc. - individual arguments by position
//	$argname - named arguments (when task.Args is defined)
func interpolateArgs(command string, args []string, task *models.Task) string {
	result := command

	// Replace $@ with all arguments
	if strings.Contains(result, "$@") {
		allArgs := strings.Join(args, " ")
		result = strings.ReplaceAll(result, "$@", allArgs)
	}

	// Replace numbered arguments $1, $2, etc.
	for i, arg := range args {
		placeholder := "$" + strconv.Itoa(i+1)
		if strings.Contains(result, placeholder) {
			result = strings.ReplaceAll(result, placeholder, arg)
		}
	}

	// Replace named arguments if task.Args is defined
	if task.Args != nil {
		allArgNames := append(task.Args.Required, task.Args.Optional...)
		for i, argName := range allArgNames {
			if i < len(args) {
				placeholder := "$" + argName
				if strings.Contains(result, placeholder) {
					result = strings.ReplaceAll(result, placeholder, args[i])
				}
			}
		}
	}

	return result
}

func (e *Executor) ExecuteTask(taskName string, task *models.Task, beforeHooks, afterHooks func([]string) error, updateCache func() error) error {
	if !e.DryRun && len(task.BeforeHooks) > 0 {
		if err := beforeHooks(task.BeforeHooks); err != nil {
			return err
		}
	}

	if !task.Silent {
		e.log.Task("Running task %q...", taskName)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}
	workingDir := task.WorkingDir
	if workingDir != "" {
		if err := os.Chdir(workingDir); err != nil {
			return fmt.Errorf("failed to change directory to %q: %v", workingDir, err)
		}
		defer os.Chdir(originalDir)
	}

	shell, shellArgs := e.shell.GetShellCommand()
	// Interpolate extra arguments using placeholders ($@, $1, $2, $argname, etc.)
	commandStr := interpolateArgs(task.Command, task.ExtraArgs, task)
	cmdArgs := append(shellArgs, commandStr)
	cmd := exec.Command(shell, cmdArgs...)

	cmd.Env = os.Environ()
	for key, value := range task.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	if task.Silent {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.Stdin = os.Stdin

	var cmdErr error
	if task.Timeout != "" {
		timeout, err := time.ParseDuration(task.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout format %q: %v", task.Timeout, err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		cmd = exec.CommandContext(ctx, shell, cmdArgs...)
		cmd.Env = os.Environ()
		for key, value := range task.Env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		}
		if task.Silent {
			cmd.Stdout = io.Discard
			cmd.Stderr = io.Discard
		} else {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		cmd.Stdin = os.Stdin

		cmdErr = cmd.Run()
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("task %q timed out after %s", taskName, task.Timeout)
		}
	} else {
		cmdErr = cmd.Run()
	}

	if cmdErr != nil {
		return fmt.Errorf("failed to run task %q: %v", taskName, cmdErr)
	}

	if err := updateCache(); err != nil {
		if !task.Silent {
			e.log.Warning("failed to update cache for task %q: %v", taskName, err)
		}
	}

	if !task.Silent {
		e.log.Success("Task %q completed successfully", taskName)
	}

	if !e.DryRun && len(task.AfterHooks) > 0 {
		if err := afterHooks(task.AfterHooks); err != nil {
			return err
		}
	}

	return nil
}

func (e *Executor) ExecuteHook(hookName string, hook *models.Hook) error {
	e.log.Task("Running hook %q...", hookName)
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}
	workingDir := hook.WorkingDir
	if workingDir != "" {
		if err := os.Chdir(workingDir); err != nil {
			return fmt.Errorf("failed to change directory to %q: %v", workingDir, err)
		}
		defer os.Chdir(originalDir)
	}

	shell, shellArgs := e.shell.GetShellCommand()
	cmdArgs := append(shellArgs, hook.Command)
	cmd := exec.Command(shell, cmdArgs...)

	cmd.Env = os.Environ()
	for key, value := range hook.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run hook %q: %v", hookName, err)
	}

	e.log.Success("Hook %q completed successfully", hookName)
	return nil
}
