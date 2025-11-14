package runner

import (
	"fmt"
	"path/filepath"
	"time"

	"azuyamat.dev/pace/internal/logger"
	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	runner   *Runner
	taskName string
	patterns []string
	log      *logger.Logger
}

func NewWatcher(runner *Runner, taskName string, patterns []string) *Watcher {
	return &Watcher{
		runner:   runner,
		taskName: taskName,
		patterns: patterns,
		log:      logger.New(),
	}
}

func (w *Watcher) Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Add files matching patterns
	for _, pattern := range w.patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			w.log.Warning("invalid pattern %q: %v", pattern, err)
			continue
		}

		for _, match := range matches {
			if err := watcher.Add(match); err != nil {
				w.log.Warning("failed to watch %q: %v", match, err)
			} else {
				w.log.Info("Watching: %s", match)
			}
		}
	}

	w.log.Info("\nWatching for changes... (Press Ctrl+C to stop)\n")

	// Run once initially
	if err := w.runner.RunTask(w.taskName); err != nil {
		w.log.Error("%v", err)
	}

	debounce := time.NewTimer(500 * time.Millisecond)
	debounce.Stop()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				w.log.Info("\nFile changed: %s", event.Name)

				// Debounce: wait for multiple rapid changes to settle
				debounce.Reset(500 * time.Millisecond)
			}

		case <-debounce.C:
			w.log.Info("\nRerunning task...")
			if err := w.runner.RunTask(w.taskName); err != nil {
				w.log.Error("%v", err)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			w.log.Error("Watcher error: %v", err)
		}
	}
}
