package runner

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"azuyamat.dev/pace/internal/config"
)

type Runner struct {
	Config *config.Config
}

func NewRunner(cfg *config.Config) *Runner {
	return &Runner{
		Config: cfg,
	}
}

func (r *Runner) RunTask(taskName string) error {
	task, exists := r.Config.Tasks[taskName]
	if !exists {
		return nil
	}

	dependencies := task.Dependencies
	for _, dependency := range dependencies {
		if err := r.RunTask(dependency); err != nil {
			return err
		}
	}

	needsRun, err := r.needsRerun(taskName)
	if err != nil {
		return fmt.Errorf("failed to check cache for task %q: %v", taskName, err)
	}

	if !needsRun {
		fmt.Printf("Task %q is up to date (cache hit)\n", taskName)
		return nil
	}

	fmt.Printf("Running task %q...\n", taskName)
	command := strings.Split(task.Command, " ")
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

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		command = append([]string{"-Command"}, command...)
		cmd = exec.Command("powershell.exe", command...)
	} else {
		command = append([]string{"-c"}, command...)
		cmd = exec.Command("sh", command...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run task %q: %v", taskName, err)
	}

	if err := r.updateCache(taskName); err != nil {
		fmt.Printf("Warning: failed to update cache for task %q: %v\n", taskName, err)
	}

	fmt.Printf("\nTask %q completed successfully.\n", taskName)
	return nil
}
