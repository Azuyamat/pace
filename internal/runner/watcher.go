package runner

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/azuyamat/pace/internal/logger"
	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	runner     *Runner
	taskName   string
	patterns   []string
	log        *logger.Logger
	cancelFunc context.CancelFunc
	taskMu     sync.Mutex
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

	if err := w.setupWatchPaths(watcher); err != nil {
		return err
	}

	w.log.Info("\nWatching for changes... (Press Ctrl+C to stop)\n")

	w.runTask()

	return w.eventLoop(watcher)
}

func (w *Watcher) setupWatchPaths(watcher *fsnotify.Watcher) error {
	dirs := make(map[string]bool)

	for _, pattern := range w.patterns {
		matches, err := w.expandPattern(pattern)
		if err != nil {
			w.log.Warning("invalid pattern %q: %v", pattern, err)
			continue
		}

		w.log.Debug("Pattern %q matched %d files", pattern, len(matches))
		for _, match := range matches {
			w.log.Debug("  - %s", match)
		}

		for _, match := range matches {
			dir := filepath.Dir(match)
			if !dirs[dir] {
				dirs[dir] = true
				if err := watcher.Add(dir); err != nil {
					w.log.Warning("failed to watch directory %q: %v", dir, err)
				} else {
					w.log.Info("Watching directory: %s", dir)
				}
			}
		}
	}

	if len(dirs) == 0 {
		return fmt.Errorf("no valid paths to watch")
	}

	return nil
}

func (w *Watcher) expandPattern(pattern string) ([]string, error) {
	if !containsDoublestar(pattern) {
		return filepath.Glob(pattern)
	}

	return w.expandDoublestar(pattern)
}

func containsDoublestar(pattern string) bool {
	return len(pattern) >= 2 && (pattern[:2] == "**" ||
		(len(pattern) >= 3 && pattern[len(pattern)-3:] == "/**") ||
		strings.Contains(pattern, "/**/") || strings.Contains(pattern, "\\**\\"))
}

func (w *Watcher) expandDoublestar(pattern string) ([]string, error) {
	var matches []string

	parts := splitPattern(pattern)
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid pattern")
	}

	baseDir := "."
	filePattern := parts[len(parts)-1]

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() {
			if matched, _ := filepath.Match(filePattern, info.Name()); matched {
				matches = append(matches, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return matches, nil
}

func splitPattern(pattern string) []string {
	var parts []string
	var current strings.Builder

	for i := 0; i < len(pattern); i++ {
		if pattern[i] == '/' || pattern[i] == '\\' {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		} else {
			current.WriteByte(pattern[i])
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

func (w *Watcher) eventLoop(watcher *fsnotify.Watcher) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	debounce := time.NewTimer(0)
	if !debounce.Stop() {
		<-debounce.C
	}

	taskDone := make(chan struct{}, 1)

	for {
		select {
		case <-sigChan:
			w.log.Info("\nShutting down watcher...")
			w.cancelCurrentTask()
			return nil

		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			w.handleEvent(event, debounce)

		case <-debounce.C:
			w.cancelCurrentTask()
			w.log.Info("\nRerunning task...")
			go w.runTaskAsync(taskDone)

		case <-taskDone:

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			w.log.Error("Watcher error: %v", err)
		}
	}
}

func (w *Watcher) handleEvent(event fsnotify.Event, debounce *time.Timer) {
	if !w.isRelevantEvent(event) {
		return
	}

	if !w.matchesPattern(event.Name) {
		return
	}

	w.log.Info("\nFile changed: %s (%s)", event.Name, event.Op.String())
	w.resetDebounce(debounce)
}

func (w *Watcher) isRelevantEvent(event fsnotify.Event) bool {
	return event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0
}

func (w *Watcher) matchesPattern(filePath string) bool {
	for _, pattern := range w.patterns {
		if containsDoublestar(pattern) {
			filePattern := filepath.Base(pattern)
			if matched, err := filepath.Match(filePattern, filepath.Base(filePath)); err == nil && matched {
				return true
			}
			continue
		}

		matched, err := filepath.Match(pattern, filePath)
		if err == nil && matched {
			return true
		}

		normalizedPattern := filepath.FromSlash(pattern)
		if matched, err := filepath.Match(normalizedPattern, filePath); err == nil && matched {
			return true
		}

		if matched, err := filepath.Match(filepath.Base(pattern), filepath.Base(filePath)); err == nil && matched {
			patternDir := filepath.Dir(pattern)
			fileDir := filepath.Dir(filePath)
			if patternDir == fileDir || patternDir == "." {
				return true
			}
		}
	}
	return false
}

func (w *Watcher) resetDebounce(debounce *time.Timer) {
	if !debounce.Stop() {
		select {
		case <-debounce.C:
		default:
		}
	}
	debounce.Reset(500 * time.Millisecond)
}

func (w *Watcher) runTask() {
	if err := w.runner.RunTask(w.taskName); err != nil {
		w.log.Error("%v", err)
	}
}

func (w *Watcher) runTaskAsync(done chan<- struct{}) {
	w.taskMu.Lock()
	ctx, cancel := context.WithCancel(context.Background())
	w.cancelFunc = cancel
	w.taskMu.Unlock()

	defer func() {
		w.taskMu.Lock()
		w.cancelFunc = nil
		w.taskMu.Unlock()
		select {
		case done <- struct{}{}:
		default:
		}
	}()

	w.runner.Reset()
	w.log.Debug("Starting task in goroutine...")
	errChan := make(chan error, 1)
	go func() {
		err := w.runner.RunTask(w.taskName)
		w.log.Debug("Task goroutine finished with err: %v", err)
		errChan <- err
	}()

	select {
	case <-ctx.Done():
		w.log.Warning("Task cancelled")
	case err := <-errChan:
		w.log.Debug("Received task completion")
		if err != nil {
			w.log.Error("%v", err)
		}
	}
}

func (w *Watcher) cancelCurrentTask() {
	w.taskMu.Lock()
	defer w.taskMu.Unlock()
	if w.cancelFunc != nil {
		w.cancelFunc()
	}
}
