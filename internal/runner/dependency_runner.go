package runner

import (
	"fmt"
	"sync"

	"github.com/azuyamat/pace/internal/models"
)

type DependencyRunner struct {
	runTask func(string, ...string) error
	log     taskLogger
}

func NewDependencyRunner(runTask func(string, ...string) error, log taskLogger) *DependencyRunner {
	return &DependencyRunner{
		runTask: runTask,
		log:     log,
	}
}

func (dr *DependencyRunner) RunDependencies(task *models.Task, dependencies []string) error {
	if !task.Parallel {
		for _, dep := range dependencies {
			if err := dr.runTask(dep); err != nil {
				if !task.ContinueOnError {
					return err
				}
			}
		}
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(dependencies))

	for _, dep := range dependencies {
		wg.Add(1)
		go func(depName string) {
			defer wg.Done()
			if err := dr.runTask(depName); err != nil {
				if !task.ContinueOnError {
					errChan <- err
				}
			}
		}(dep)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

type HookExecutor struct {
	hooks    map[string]models.Hook
	executor *Executor
	log      taskLogger
}

func NewHookExecutor(hooks map[string]models.Hook, executor *Executor, log taskLogger) *HookExecutor {
	return &HookExecutor{
		hooks:    hooks,
		executor: executor,
		log:      log,
	}
}

func (he *HookExecutor) ExecuteHooks(hookNames []string) error {
	for _, hookName := range hookNames {
		hook, exists := he.hooks[hookName]
		if !exists {
			return fmt.Errorf("hook %q not found", hookName)
		}
		if err := he.executor.ExecuteHook(hookName, &hook); err != nil {
			return fmt.Errorf("hook %q failed: %v", hookName, err)
		}
	}
	return nil
}
