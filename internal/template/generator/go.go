package generator

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/azuyamat/pace/internal/config"
	"github.com/azuyamat/pace/internal/logger"
	"github.com/azuyamat/pace/internal/models"
)

func NewGoGenerator() Generator {
	return newGenerator(generate, models.ProjectTypeGo)
}

func generate() (config.Config, error) {
	cfg := config.NewDefaultConfig()

	moduleName := detectModuleName()
	if moduleName != "" {
		logger.Info("Module: %s", moduleName)
	}

	projectStructure := detectGoProjectStructure()
	logger.Info("Project structure: %s", projectStructure)

	mainPath := detectMainPackage()
	if mainPath != "" {
		logger.Info("Main package: %s", mainPath)
	}

	outputBinary := "bin/" + filepath.Base(moduleName)
	if outputBinary == "bin/" {
		outputBinary = "bin/app"
	}

	cfg.DefaultTask = "build"

	buildCmd := "go build -o " + outputBinary
	if mainPath != "" {
		buildCmd += " ./" + mainPath
	} else {
		buildCmd += " ./..."
	}

	cfg.Tasks["build"] = models.Task{
		Name:        "build",
		Alias:       "b",
		Command:     buildCmd,
		Description: "Build the application",
		Cache:       true,
		Inputs:      []string{"**/*.go", "go.mod", "go.sum"},
		Outputs:     []string{outputBinary},
	}

	cfg.Tasks["test"] = models.Task{
		Name:        "test",
		Alias:       "t",
		Command:     "go test ./...",
		Description: "Run tests",
		Cache:       true,
		Inputs:      []string{"**/*.go"},
	}

	if hasGoLintConfig() {
		cfg.Tasks["lint"] = models.Task{
			Name:        "lint",
			Alias:       "l",
			Command:     "golangci-lint run",
			Description: "Lint code",
			Cache:       true,
			Inputs:      []string{"**/*.go", ".golangci.yml", ".golangci.yaml"},
		}
	}

	runCmd := "go run"
	if mainPath != "" {
		runCmd += " ./" + mainPath
	} else {
		runCmd += " ./..."
	}

	cfg.Tasks["run"] = models.Task{
		Name:        "run",
		Alias:       "r",
		Command:     runCmd,
		Description: "Run the application",
		Watch:       true,
		Inputs:      []string{"**/*.go"},
	}

	if hasToolsGo() {
		cfg.Tasks["tools"] = models.Task{
			Name:        "tools",
			Command:     "go generate -tags tools ./tools",
			Description: "Install development tools",
			Inputs:      []string{"tools/tools.go"},
		}
	}

	if hasGoGenerate() {
		cfg.Tasks["generate"] = models.Task{
			Name:        "generate",
			Command:     "go generate ./...",
			Description: "Run go generate",
			Inputs:      []string{"**/*.go"},
		}
	}

	cfg.Tasks["tidy"] = models.Task{
		Name:        "tidy",
		Command:     "go mod tidy",
		Description: "Tidy go.mod and go.sum",
		Inputs:      []string{"go.mod"},
	}

	cfg.Tasks["vet"] = models.Task{
		Name:        "vet",
		Command:     "go vet ./...",
		Description: "Run go vet",
		Cache:       true,
		Inputs:      []string{"**/*.go"},
	}

	cleanCommand := "rm -rf bin"
	if isWindows() {
		cleanCommand = "if exist bin rmdir /s /q bin"
	}

	cfg.Hooks["clean"] = models.Hook{
		Name:        "clean",
		Command:     cleanCommand,
		Description: "Clean build artifacts",
	}

	cfg.Hooks["fmt"] = models.Hook{
		Name:        "fmt",
		Command:     "go fmt ./...",
		Description: "Format code",
	}

	return *cfg, nil
}

func detectModuleName() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return ""
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}
	return ""
}

func detectGoProjectStructure() string {
	hasCmdDir := hasDirectory("cmd")
	hasPkgDir := hasDirectory("pkg")
	hasInternalDir := hasDirectory("internal")

	if hasCmdDir && (hasPkgDir || hasInternalDir) {
		return "standard"
	}
	if hasCmdDir {
		return "cmd-based"
	}
	if hasFile("main.go") {
		return "simple"
	}
	return "library"
}

func detectMainPackage() string {
	if hasFile("main.go") {
		return ""
	}

	cmdDirs := []string{"cmd/server", "cmd/main", "cmd/app"}
	for _, dir := range cmdDirs {
		if hasFile(filepath.Join(dir, "main.go")) {
			return dir
		}
	}

	if hasDirectory("cmd") {
		entries, err := os.ReadDir("cmd")
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					mainPath := filepath.Join("cmd", entry.Name(), "main.go")
					if hasFile(mainPath) {
						return filepath.Join("cmd", entry.Name())
					}
				}
			}
		}
	}

	return ""
}

func hasGoLintConfig() bool {
	return hasFile(".golangci.yml") || hasFile(".golangci.yaml") || hasFile(".golangci.toml")
}

func hasToolsGo() bool {
	return hasFile("tools/tools.go")
}

func hasGoGenerate() bool {
	return grepGoFiles("//go:generate")
}

func grepGoFiles(pattern string) bool {
	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		if strings.Contains(string(data), pattern) {
			return filepath.SkipAll
		}

		return nil
	}) == filepath.SkipAll
}
