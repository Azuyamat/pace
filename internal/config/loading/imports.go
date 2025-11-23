package loading

import (
	"fmt"
	"path/filepath"
)

func processImports(cfg *Config, baseDir string) error {
	visited := make(map[string]bool)
	return processImportsRecursive(cfg, baseDir, visited)
}

func processImportsRecursive(cfg *Config, baseDir string, visited map[string]bool) error {
	for _, importPath := range cfg.Imports {
		fullPath := filepath.Join(baseDir, importPath)
		absPath, err := filepath.Abs(fullPath)
		if err != nil {
			return fmt.Errorf("failed to resolve import path %q: %v", fullPath, err)
		}

		if visited[absPath] {
			return fmt.Errorf("circular import detected: %q", absPath)
		}

		visited[absPath] = true

		importedCfg, err := ParseFile(fullPath)
		if err != nil {
			return err
		}

		if err := processImportsRecursive(importedCfg, filepath.Dir(fullPath), visited); err != nil {
			return err
		}

		importField(importedCfg.Tasks, cfg.Tasks)
		importField(importedCfg.Hooks, cfg.Hooks)
		importField(importedCfg.Constants, cfg.Constants)
		importField(importedCfg.Globals, cfg.Globals)
		importField(importedCfg.Aliases, cfg.Aliases)
	}

	return nil
}

func importField[T any](src, dest map[string]T) {
	for name, value := range src {
		if _, exists := dest[name]; !exists {
			dest[name] = value
		}
	}
}
