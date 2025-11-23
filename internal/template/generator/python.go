package generator

import (
	"os"
	"strings"

	"github.com/azuyamat/pace/internal/config"
	"github.com/azuyamat/pace/internal/logger"
	"github.com/azuyamat/pace/internal/models"
)

func NewPythonGenerator() Generator {
	return newGenerator(generatePython, models.ProjectTypePython)
}

func generatePython() (config.Config, error) {
	cfg := config.NewDefaultConfig()

	depManager := detectPythonDependencyManager()
	logger.Info("Using dependency manager: %s", depManager)

	testFramework := detectTestFramework()
	logger.Info("Detected test framework: %s", testFramework)

	linter := detectLinter()
	formatter := detectFormatter()

	entryPoint := detectEntryPoint()
	if entryPoint != "" {
		logger.Info("Detected entry point: %s", entryPoint)
		cfg.DefaultTask = "run"
	} else {
		cfg.DefaultTask = "test"
	}

	cfg.Hooks["install"] = generateInstallHook(depManager)

	if entryPoint != "" {
		cfg.Tasks["run"] = models.Task{
			Name:        "run",
			Alias:       "r",
			Command:     "python " + entryPoint,
			Description: "Run the application",
			Requires:    []string{"install"},
			Inputs:      []string{"**/*.py"},
		}
	}

	cfg.Tasks["test"] = generateTestTask(testFramework)

	if linter != "" {
		cfg.Tasks["lint"] = generateLintTask(linter)
	}

	if formatter != "" {
		cfg.Tasks["format"] = generateFormatTask(formatter)
	}

	if hasPyprojectToml() {
		cfg.Tasks["type-check"] = models.Task{
			Name:        "type-check",
			Command:     "mypy .",
			Description: "Type check with mypy",
			Requires:    []string{"install"},
			Cache:       true,
			Inputs:      []string{"**/*.py", "pyproject.toml"},
		}
	}

	cleanCommand := "rm -rf __pycache__ .pytest_cache .coverage .mypy_cache dist build *.egg-info"
	if isWindows() {
		cleanCommand = "for /d /r . %d in (__pycache__) do @if exist \"%d\" rd /s /q \"%d\""
	}

	cfg.Hooks["clean"] = models.Hook{
		Name:        "clean",
		Command:     cleanCommand,
		Description: "Clean Python artifacts",
	}

	return *cfg, nil
}

func detectPythonDependencyManager() string {
	if _, err := os.Stat("pyproject.toml"); err == nil {
		data, err := os.ReadFile("pyproject.toml")
		if err == nil {
			content := string(data)
			if strings.Contains(content, "[tool.poetry]") {
				return "poetry"
			}
			if strings.Contains(content, "[tool.pdm]") {
				return "pdm"
			}
		}
	}
	if _, err := os.Stat("Pipfile"); err == nil {
		return "pipenv"
	}
	return "pip"
}

func detectTestFramework() string {
	if hasFile("pytest.ini") || hasFile("pyproject.toml") {
		return "pytest"
	}
	if hasDirectory("tests") || hasDirectory("test") {
		return "pytest"
	}
	return "unittest"
}

func detectLinter() string {
	if hasFile(".ruff.toml") || hasFile("ruff.toml") {
		return "ruff"
	}
	if hasFile(".flake8") || hasFile("setup.cfg") {
		return "flake8"
	}
	if hasFile("pylintrc") || hasFile(".pylintrc") {
		return "pylint"
	}
	return "ruff"
}

func detectFormatter() string {
	if hasFile(".ruff.toml") || hasFile("ruff.toml") {
		return "ruff"
	}
	if hasFile("pyproject.toml") {
		data, err := os.ReadFile("pyproject.toml")
		if err == nil && strings.Contains(string(data), "[tool.black]") {
			return "black"
		}
	}
	return "black"
}

func detectEntryPoint() string {
	entryPoints := []string{"main.py", "app.py", "run.py", "src/main.py", "src/app.py"}
	for _, entry := range entryPoints {
		if hasFile(entry) {
			return entry
		}
	}
	return ""
}

func hasPyprojectToml() bool {
	return hasFile("pyproject.toml")
}

func generateInstallHook(depManager string) models.Hook {
	var command string

	switch depManager {
	case "poetry":
		command = "poetry install"
	case "pdm":
		command = "pdm install"
	case "pipenv":
		command = "pipenv install"
	default:
		command = "pip install -r requirements.txt"
	}

	return models.Hook{
		Name:        "install",
		Command:     command,
		Description: "Install dependencies",
	}
}

func generateTestTask(framework string) models.Task {
	var command string

	switch framework {
	case "pytest":
		command = "pytest"
	case "unittest":
		command = "python -m unittest discover"
	default:
		command = "pytest"
	}

	return models.Task{
		Name:        "test",
		Alias:       "t",
		Command:     command,
		Description: "Run tests",
		Requires:    []string{"install"},
		Cache:       true,
		Inputs:      []string{"src/**/*.py", "tests/**/*.py", "test/**/*.py"},
	}
}

func generateLintTask(linter string) models.Task {
	var command string
	var inputs []string

	switch linter {
	case "ruff":
		command = "ruff check ."
		inputs = []string{"**/*.py", "ruff.toml", ".ruff.toml"}
	case "flake8":
		command = "flake8 ."
		inputs = []string{"**/*.py", ".flake8", "setup.cfg"}
	case "pylint":
		command = "pylint **/*.py"
		inputs = []string{"**/*.py", "pylintrc", ".pylintrc"}
	default:
		command = "ruff check ."
		inputs = []string{"**/*.py"}
	}

	return models.Task{
		Name:        "lint",
		Alias:       "l",
		Command:     command,
		Description: "Lint code",
		Requires:    []string{"install"},
		Cache:       true,
		Inputs:      inputs,
	}
}

func generateFormatTask(formatter string) models.Task {
	var command string
	var inputs []string

	switch formatter {
	case "ruff":
		command = "ruff format ."
		inputs = []string{"**/*.py", "ruff.toml", ".ruff.toml"}
	case "black":
		command = "black ."
		inputs = []string{"**/*.py", "pyproject.toml"}
	default:
		command = "black ."
		inputs = []string{"**/*.py"}
	}

	return models.Task{
		Name:        "format",
		Alias:       "f",
		Command:     command,
		Description: "Format code",
		Inputs:      inputs,
	}
}
