package generator

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/azuyamat/pace/internal/config"
	"github.com/azuyamat/pace/internal/logger"
	"github.com/azuyamat/pace/internal/models"
)

func NewRustGenerator() Generator {
	return newGenerator(generateRust, models.ProjectTypeRust)
}

func generateRust() (config.Config, error) {
	cfg := config.NewDefaultConfig()

	projectType := detectRustProjectType()
	logger.Info("Detected Rust project type: %s", projectType)

	packageName := detectRustPackageName()
	if packageName != "" {
		logger.Info("Package name: %s", packageName)
	}

	workspaceMembers := detectWorkspaceMembers()
	if len(workspaceMembers) > 0 {
		logger.Info("Found workspace with %d members", len(workspaceMembers))
	}

	cfg.DefaultTask = "build"

	cfg.Tasks["build"] = models.Task{
		Name:        "build",
		Alias:       "b",
		Command:     "cargo build",
		Description: "Build the project",
		Cache:       true,
		Inputs:      []string{"src/**/*.rs", "Cargo.toml", "Cargo.lock"},
		Outputs:     []string{"target/debug"},
	}

	cfg.Tasks["release"] = models.Task{
		Name:        "release",
		Alias:       "r",
		Command:     "cargo build --release",
		Description: "Build release binary",
		Requires:    []string{"test"},
		Cache:       true,
		Inputs:      []string{"src/**/*.rs", "Cargo.toml", "Cargo.lock"},
		Outputs:     []string{"target/release"},
	}

	cfg.Tasks["test"] = models.Task{
		Name:        "test",
		Alias:       "t",
		Command:     "cargo test",
		Description: "Run tests",
		Cache:       true,
		Inputs:      []string{"src/**/*.rs", "tests/**/*.rs"},
	}

	cfg.Tasks["check"] = models.Task{
		Name:        "check",
		Alias:       "c",
		Command:     "cargo check",
		Description: "Check code without building",
		Cache:       true,
		Inputs:      []string{"src/**/*.rs", "Cargo.toml"},
	}

	cfg.Tasks["lint"] = models.Task{
		Name:        "lint",
		Alias:       "l",
		Command:     "cargo clippy -- -D warnings",
		Description: "Lint code with clippy",
		Cache:       true,
		Inputs:      []string{"src/**/*.rs"},
	}

	cfg.Tasks["format"] = models.Task{
		Name:        "format",
		Alias:       "f",
		Command:     "cargo fmt",
		Description: "Format code",
		Inputs:      []string{"src/**/*.rs"},
	}

	cfg.Tasks["format-check"] = models.Task{
		Name:        "format-check",
		Command:     "cargo fmt -- --check",
		Description: "Check code formatting",
		Cache:       true,
		Inputs:      []string{"src/**/*.rs"},
	}

	if projectType == "binary" || projectType == "both" {
		cfg.Tasks["run"] = models.Task{
			Name:        "run",
			Command:     "cargo run",
			Description: "Run the application",
			DependsOn:   []string{"build"},
			Watch:       true,
			Inputs:      []string{"src/**/*.rs"},
		}
	}

	if hasBenchmarks() {
		cfg.Tasks["bench"] = models.Task{
			Name:        "bench",
			Command:     "cargo bench",
			Description: "Run benchmarks",
			Inputs:      []string{"benches/**/*.rs", "src/**/*.rs"},
		}
	}

	if hasExamples() {
		cfg.Tasks["examples"] = models.Task{
			Name:        "examples",
			Command:     "cargo build --examples",
			Description: "Build all examples",
			Cache:       true,
			Inputs:      []string{"examples/**/*.rs", "src/**/*.rs"},
		}
	}

	cfg.Hooks["clean"] = models.Hook{
		Name:        "clean",
		Command:     "cargo clean",
		Description: "Clean build artifacts",
	}

	return *cfg, nil
}

func detectRustProjectType() string {
	data, err := os.ReadFile("Cargo.toml")
	if err != nil {
		return "unknown"
	}

	content := string(data)
	hasBin := strings.Contains(content, "[[bin]]") || hasFile("src/main.rs")
	hasLib := strings.Contains(content, "[lib]") || hasFile("src/lib.rs")

	if hasBin && hasLib {
		return "both"
	}
	if hasBin {
		return "binary"
	}
	if hasLib {
		return "library"
	}
	return "unknown"
}

func detectRustPackageName() string {
	data, err := os.ReadFile("Cargo.toml")
	if err != nil {
		return ""
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "name") && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[1])
				name = strings.Trim(name, "\"")
				return name
			}
		}
	}
	return ""
}

func detectWorkspaceMembers() []string {
	data, err := os.ReadFile("Cargo.toml")
	if err != nil {
		return nil
	}

	content := string(data)
	if !strings.Contains(content, "[workspace]") {
		return nil
	}

	var members []string
	lines := strings.Split(content, "\n")
	inWorkspace := false
	inMembers := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "[workspace]" {
			inWorkspace = true
			continue
		}

		if inWorkspace && strings.HasPrefix(line, "members") {
			inMembers = true
			continue
		}

		if inMembers && strings.HasPrefix(line, "[") {
			break
		}

		if inMembers && strings.Contains(line, "\"") {
			member := strings.Trim(line, " \t,[]\"")
			if member != "" {
				members = append(members, member)
			}
		}
	}

	return members
}

func hasBenchmarks() bool {
	return hasDirectory("benches")
}

func hasExamples() bool {
	if !hasDirectory("examples") {
		return false
	}

	entries, err := os.ReadDir("examples")
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".rs" {
			return true
		}
	}

	return false
}
