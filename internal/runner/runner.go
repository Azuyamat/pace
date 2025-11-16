package runner

import (
	"fmt"
	"sync"
	"time"

	"azuyamat.dev/pace/internal/config"
	"azuyamat.dev/pace/internal/logger"
	"azuyamat.dev/pace/internal/models"
)

type Runner struct {
	Config           *config.Config
	completed        map[string]bool
	running          map[string]bool
	mu               sync.Mutex
	DryRun           bool
	Force            bool
	log              *logger.Logger
	shell            *Shell
	executor         *Executor
	dependencyRunner *DependencyRunner
	hookExecutor     *HookExecutor
}

func NewRunner(cfg *config.Config) *Runner {
	log := logger.New()
	shell := NewShell(cfg.Globals)
	executor := NewExecutor(shell, log, false)

	r := &Runner{
		Config:    cfg,
		completed: make(map[string]bool),
		running:   make(map[string]bool),
		log:       log,
		shell:     shell,
		executor:  executor,
	}

	r.dependencyRunner = NewDependencyRunner(r.RunTask, log)
	r.hookExecutor = NewHookExecutor(cfg.Hooks, executor, log)

	return r
}

func (r *Runner) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.completed = make(map[string]bool)
	r.running = make(map[string]bool)
}

func (r *Runner) validateAndSetArgs(task *models.Task, extraArgs []string) error {
	// If no args definition, use old behavior (positional only)
	if task.Args == nil {
		task.ExtraArgs = extraArgs
		return nil
	}

	// Check if we have enough required arguments
	requiredCount := len(task.Args.Required)
	if len(extraArgs) < requiredCount {
		return fmt.Errorf("task %q requires %d argument(s) but got %d\nRequired: %v",
			task.Name, requiredCount, len(extraArgs), task.Args.Required)
	}

	// Check if we have too many arguments
	totalExpected := len(task.Args.Required) + len(task.Args.Optional)
	if len(extraArgs) > totalExpected {
		return fmt.Errorf("task %q expects at most %d argument(s) but got %d\nExpected: required=%v, optional=%v",
			task.Name, totalExpected, len(extraArgs), task.Args.Required, task.Args.Optional)
	}

	task.ExtraArgs = extraArgs
	return nil
}

func (r *Runner) RunTask(taskName string, extraArgs ...string) error {
	task, exists := r.Config.Tasks[taskName]
	if !exists {
		return fmt.Errorf("task %q not found", taskName)
	}

	if err := r.validateAndSetArgs(&task, extraArgs); err != nil {
		return err
	}

	r.mu.Lock()
	if r.completed[taskName] {
		r.mu.Unlock()
		return nil
	}
	if r.running[taskName] {
		r.mu.Unlock()
		return fmt.Errorf("circular dependency detected for task %q", taskName)
	}
	r.running[taskName] = true
	r.mu.Unlock()

	if len(task.Dependencies) > 0 {
		if err := r.dependencyRunner.RunDependencies(&task, task.Dependencies); err != nil {
			return err
		}
	}

	needsRun := true
	if r.Force {
		needsRun = true
	} else {
		var err error
		needsRun, err = r.needsRerun(taskName)
		if err != nil {
			return fmt.Errorf("failed to check cache for task %q: %v", taskName, err)
		}

		if !needsRun {
			if !task.Silent {
				r.log.Info("Task %q is up to date (cache hit)", taskName)
			}
			r.mu.Lock()
			r.completed[taskName] = true
			r.running[taskName] = false
			r.mu.Unlock()
			return nil
		}
	}

	if r.DryRun {
		cmdStr := interpolateArgs(task.Command, task.ExtraArgs, &task)
		if len(task.ExtraArgs) > 0 && cmdStr == task.Command {
			// Arguments provided but not used in command
			r.log.Warning("[DRY RUN] Extra arguments provided but command has no placeholders ($@, $1, $2, etc.): %v", task.ExtraArgs)
		}
		r.log.Debug("[DRY RUN] Would execute task %q: %s", taskName, cmdStr)
		if len(task.BeforeHooks) > 0 {
			r.log.Debug("[DRY RUN] Would run before hooks: %v", task.BeforeHooks)
		}
		if len(task.AfterHooks) > 0 {
			r.log.Debug("[DRY RUN] Would run after hooks: %v", task.AfterHooks)
		}
		if len(task.OnSuccess) > 0 {
			r.log.Debug("[DRY RUN] Would run on_success hooks: %v", task.OnSuccess)
		}
		r.mu.Lock()
		r.completed[taskName] = true
		r.running[taskName] = false
		r.mu.Unlock()
		return nil
	}

	r.executor.DryRun = r.DryRun

	var execErr error
	for attempt := 0; attempt <= task.Retry; attempt++ {
		if attempt > 0 {
			if !task.Silent {
				r.log.Warning("Retrying task %q (attempt %d/%d)...", taskName, attempt, task.Retry)
			}
			if task.RetryDelay != "" {
				delay, err := time.ParseDuration(task.RetryDelay)
				if err == nil {
					time.Sleep(delay)
				}
			}
		}

		beforeHookFunc := func(hooks []string) error {
			return r.hookExecutor.ExecuteHooks(hooks)
		}
		afterHookFunc := func(hooks []string) error {
			return r.hookExecutor.ExecuteHooks(hooks)
		}
		updateCacheFunc := func() error {
			return r.updateCache(taskName)
		}

		execErr = r.executor.ExecuteTask(taskName, &task, beforeHookFunc, afterHookFunc, updateCacheFunc)
		if execErr == nil {
			break
		}
	}

	if execErr != nil {
		if !r.DryRun && len(task.OnFailure) > 0 {
			if err := r.hookExecutor.ExecuteHooks(task.OnFailure); err != nil {
				if !task.Silent {
					r.log.Warning("failure hook execution failed: %v", err)
				}
			}
		}
		return execErr
	}

	if !r.DryRun && len(task.OnSuccess) > 0 {
		if err := r.hookExecutor.ExecuteHooks(task.OnSuccess); err != nil {
			if !task.Silent {
				r.log.Warning("success hook execution failed: %v", err)
			}
		}
	}

	r.mu.Lock()
	r.completed[taskName] = true
	r.running[taskName] = false
	r.mu.Unlock()

	return nil
}
