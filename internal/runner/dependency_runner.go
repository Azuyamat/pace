package runner

import (
	"context"
	"fmt"
	"sync"

	"github.com/azuyamat/pace/internal/models"
)

type DependencyRunner struct {
	runTask         func(task models.Task, extraArgs ...string) error
	runTaskWithCtx  func(ctx context.Context, task models.Task, extraArgs ...string) error
	log             taskLogger
}

func NewDependencyRunner(runTask func(models.Task, ...string) error, log taskLogger) *DependencyRunner {
	return &DependencyRunner{
		runTask: runTask,
		log:     log,
	}
}

func (dr *DependencyRunner) SetContextRunner(runTaskWithCtx func(context.Context, models.Task, ...string) error) {
	dr.runTaskWithCtx = runTaskWithCtx
}

func (dr *DependencyRunner) RunDependencies(task *models.Task, dependencies []models.Task) error {
	return dr.RunDependenciesWithContext(context.Background(), task, dependencies)
}

func (dr *DependencyRunner) RunDependenciesWithContext(ctx context.Context, task *models.Task, dependencies []models.Task) error {
	runFunc := dr.runTask
	if dr.runTaskWithCtx != nil {
		runFunc = func(dep models.Task, extraArgs ...string) error {
			return dr.runTaskWithCtx(ctx, dep, extraArgs...)
		}
	}

	if !task.Parallel {
		for _, dep := range dependencies {
			if err := runFunc(dep); err != nil {
				if !task.ContinueOnError {
					return err
				}
			}
		}
		return nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	for _, dep := range dependencies {
		wg.Add(1)
		go func(dep models.Task) {
			defer wg.Done()

			if err := runFunc(dep); err != nil {
				if !task.ContinueOnError {
					mu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					mu.Unlock()
				}
			}
		}(dep)
	}

	wg.Wait()

	return firstErr
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
