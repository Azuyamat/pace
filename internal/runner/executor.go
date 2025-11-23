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

	"github.com/azuyamat/pace/internal/models"
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
	replaced := make(map[string]bool)

	if task.Args != nil {
		allArgNames := append(task.Args.Required, task.Args.Optional...)
		for i, argName := range allArgNames {
			if i < len(args) {
				placeholder := "$" + argName
				key := placeholder + "_named"
				if strings.Contains(result, placeholder) && !replaced[key] {
					result = strings.ReplaceAll(result, placeholder, args[i])
					replaced[key] = true
				}
			}
		}
	}

	for i := len(args); i >= 1; i-- {
		placeholder := "$" + strconv.Itoa(i)
		key := placeholder
		if strings.Contains(result, placeholder) && !replaced[key] {
			result = strings.ReplaceAll(result, placeholder, args[i-1])
			replaced[key] = true
		}
	}

	if strings.Contains(result, "$@") && !replaced["$@"] {
		allArgs := strings.Join(args, " ")
		result = strings.ReplaceAll(result, "$@", allArgs)
	}

	return result
}

func (e *Executor) ExecuteTask(taskName string, task *models.Task, beforeHooks, afterHooks func([]string) error, updateCache func() error) error {
	return e.ExecuteTaskWithContext(context.Background(), taskName, task, beforeHooks, afterHooks, updateCache)
}

func (e *Executor) ExecuteTaskWithContext(ctx context.Context, taskName string, task *models.Task, beforeHooks, afterHooks func([]string) error, updateCache func() error) error {
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
		defer func() {
			if err := os.Chdir(originalDir); err != nil {
				e.log.Warning("failed to revert to original directory %q: %v", originalDir, err)
			}
		}()
	}

	shell, shellArgs := e.shell.GetShellCommand()
	commandStr := interpolateArgs(task.Command, task.ExtraArgs, task)
	cmdArgs := append(shellArgs, commandStr)

	execCtx := ctx
	if task.Timeout != "" {
		timeout, err := time.ParseDuration(task.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout format %q: %v", task.Timeout, err)
		}
		var cancel context.CancelFunc
		execCtx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(execCtx, shell, cmdArgs...)
	cmd.Env = os.Environ()
	for key, value := range task.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	if task.Silent {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	} else {
		stdoutWriter := NewPrefixedWriter(taskName, true)
		stderrWriter := NewPrefixedWriter(taskName, false)
		defer stdoutWriter.Close()
		defer stderrWriter.Close()

		cmd.Stdout = stdoutWriter
		cmd.Stderr = stderrWriter
	}
	cmd.Stdin = os.Stdin

	cmdErr := cmd.Run()

	if execCtx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("task %q timed out after %s", taskName, task.Timeout)
	}

	if execCtx.Err() == context.Canceled {
		return fmt.Errorf("task %q was cancelled", taskName)
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

	stdoutWriter := NewPrefixedWriter(hookName, true)
	stderrWriter := NewPrefixedWriter(hookName, false)
	defer stdoutWriter.Close()
	defer stderrWriter.Close()

	cmd.Stdout = stdoutWriter
	cmd.Stderr = stderrWriter
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run hook %q: %v", hookName, err)
	}

	e.log.Success("Hook %q completed successfully", hookName)
	return nil
}
