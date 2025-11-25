package detector

import "os"

type indicatorFunc func(projectPath string) bool

func fileExistsIndicator(filePath string) indicatorFunc {
	return func(projectPath string) bool {
		fullPath := projectPath + "/" + filePath
		if info, err := os.Stat(fullPath); os.IsNotExist(err) {
			return false
		} else if err == nil {
			return !info.IsDir()
		}
		return false
	}
}

func dirExistsIndicator(dirPath string) indicatorFunc {
	return func(projectPath string) bool {
		fullPath := projectPath + "/" + dirPath
		if info, err := os.Stat(fullPath); os.IsNotExist(err) {
			return false
		} else if err == nil {
			return info.IsDir()
		}
		return false
	}
}

func anyFileExistsIndicator(filePaths ...string) indicatorFunc {
	return func(projectPath string) bool {
		for _, filePath := range filePaths {
			if fileExistsIndicator(filePath)(projectPath) {
				return true
			}
		}
		return false
	}
}
