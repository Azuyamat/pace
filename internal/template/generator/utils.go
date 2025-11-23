package generator

import (
	"os"
)

func hasFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func hasDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
