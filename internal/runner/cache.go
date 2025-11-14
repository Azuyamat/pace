package runner

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const cacheDir = ".pace-cache"

type TaskCache struct {
	TaskName     string    `json:"task_name"`
	InputsHash   string    `json:"inputs_hash"`
	OutputsHash  string    `json:"outputs_hash"`
	LastRunTime  time.Time `json:"last_run_time"`
	CommandHash  string    `json:"command_hash"`
	Dependencies []string  `json:"dependencies"`
}

func computeFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func computeFilesHash(patterns []string) (string, error) {
	if len(patterns) == 0 {
		return "", nil
	}

	hash := sha256.New()
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return "", fmt.Errorf("invalid pattern %q: %v", pattern, err)
		}

		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				continue // File doesn't exist yet, skip
			}
			if info.IsDir() {
				continue // Skip directories
			}

			fileHash, err := computeFileHash(match)
			if err != nil {
				return "", err
			}
			hash.Write([]byte(match + ":" + fileHash))
		}
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func computeStringHash(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func ensureCacheDir() error {
	return os.MkdirAll(cacheDir, 0755)
}

func getCachePath(taskName string) string {
	return filepath.Join(cacheDir, taskName+".json")
}

func loadCache(taskName string) (*TaskCache, error) {
	cachePath := getCachePath(taskName)

	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No cache exists
		}
		return nil, err
	}

	var cache TaskCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

func saveCache(cache *TaskCache) error {
	if err := ensureCacheDir(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	cachePath := getCachePath(cache.TaskName)
	return os.WriteFile(cachePath, data, 0644)
}

func (r *Runner) needsRerun(taskName string) (bool, error) {
	task, exists := r.Config.Tasks[taskName]
	if !exists {
		return false, fmt.Errorf("task %q not found", taskName)
	}

	if !task.Cache {
		return true, nil
	}

	cache, err := loadCache(taskName)
	if err != nil {
		return false, err
	}
	if cache == nil {
		return true, nil // No cache exists, need to run
	}

	currentCommandHash := computeStringHash(task.Command)
	if cache.CommandHash != currentCommandHash {
		return true, nil // Command changed, need to run
	}

	if len(cache.Dependencies) != len(task.Dependencies) {
		return true, nil // Dependencies changed, need to run
	}
	for i, dep := range task.Dependencies {
		if i >= len(cache.Dependencies) || cache.Dependencies[i] != dep {
			return true, nil // Dependencies changed, need to run
		}
	}

	currentInputsHash, err := computeFilesHash(task.Inputs)
	if err != nil {
		return false, err
	}
	if cache.InputsHash != currentInputsHash {
		return true, nil // Inputs changed, need to run
	}

	if len(task.Outputs) > 0 {
		for _, outputPattern := range task.Outputs {
			matches, err := filepath.Glob(outputPattern)
			if err != nil {
				return false, fmt.Errorf("invalid output pattern %q: %v", outputPattern, err)
			}

			if len(matches) == 0 {
				return true, nil // Output doesn't exist, need to run
			}

			for _, match := range matches {
				info, err := os.Stat(match)
				if err != nil {
					return true, nil // Output doesn't exist, need to run
				}
				if info.IsDir() {
					continue
				}

				if info.ModTime().After(cache.LastRunTime) {
					currentOutputsHash, err := computeFilesHash(task.Outputs)
					if err != nil {
						return false, err
					}
					if cache.OutputsHash != currentOutputsHash {
						return true, nil // Outputs changed, need to run
					}
				}
			}
		}
	}

	return false, nil // Cache is valid, no need to run
}

func (r *Runner) updateCache(taskName string) error {
	task, exists := r.Config.Tasks[taskName]
	if !exists {
		return fmt.Errorf("task %q not found", taskName)
	}

	if !task.Cache {
		return nil
	}

	inputsHash, err := computeFilesHash(task.Inputs)
	if err != nil {
		return err
	}

	outputsHash, err := computeFilesHash(task.Outputs)
	if err != nil {
		return err
	}

	cache := &TaskCache{
		TaskName:     taskName,
		InputsHash:   inputsHash,
		OutputsHash:  outputsHash,
		LastRunTime:  time.Now(),
		CommandHash:  computeStringHash(task.Command),
		Dependencies: task.Dependencies,
	}

	return saveCache(cache)
}
