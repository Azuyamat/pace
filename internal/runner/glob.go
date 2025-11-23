package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func expandGlobPattern(pattern string) ([]string, error) {
	if !containsDoublestar(pattern) {
		return filepath.Glob(pattern)
	}

	return expandDoublestar(pattern)
}

func containsDoublestar(pattern string) bool {
	return len(pattern) >= 2 && (pattern[:2] == "**" ||
		(len(pattern) >= 3 && pattern[len(pattern)-3:] == "/**") ||
		strings.Contains(pattern, "/**/") || strings.Contains(pattern, "\\**\\"))
}

func expandDoublestar(pattern string) ([]string, error) {
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
