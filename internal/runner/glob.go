package runner

import (
	"io/fs"
	"path/filepath"

	"github.com/azuyamat/globber/glob"
)

func expandGlobPattern(pattern string) ([]string, error) {
	var matches []string

	fsMatcher := glob.FSMatcher(pattern)
	err := fsMatcher.WalkDirFS(".", func(path string, entry fs.DirEntry) error {
		matches = append(matches, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return matches, nil
}

func matchesGlobPattern(pattern, filePath string) bool {
	matcher := glob.Matcher(pattern)
	normalizedPath := filepath.ToSlash(filePath)
	matches, _ := matcher.Matches(normalizedPath)
	return matches
}
