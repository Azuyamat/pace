package generator

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/azuyamat/pace/internal/config"
	"github.com/azuyamat/pace/internal/logger"
	"github.com/azuyamat/pace/internal/models"
)

type packageJSON struct {
	Scripts map[string]string `json:"scripts"`
	Name    string            `json:"name"`
}

func NewNodeGenerator() Generator {
	return newGenerator(generateNode, models.ProjectTypeNode)
}

func generateNode() (config.Config, error) {
	cfg := config.NewDefaultConfig()

	pkg, err := readPackageJSON()
	if err != nil {
		logger.Warning("Could not read package.json: %v", err)
		return generateDefaultNode(cfg)
	}

	logger.Info("Found %d scripts in package.json", len(pkg.Scripts))

	packageManager := detectPackageManager()
	logger.Info("Using package manager: %s", packageManager)

	cfg.DefaultTask = determineDefaultTask(pkg.Scripts)

	cfg.Hooks["install"] = models.Hook{
		Name:        "install",
		Command:     packageManager + " install",
		Description: "Install dependencies",
	}

	for scriptName, scriptCommand := range pkg.Scripts {
		task := generateTaskFromScript(scriptName, scriptCommand, packageManager)
		if task.Name != "" {
			cfg.Tasks[task.Name] = task
		}
	}

	if _, hasBuild := cfg.Tasks["build"]; !hasBuild && pkg.Scripts["build"] != "" {
		cfg.Tasks["build"] = models.Task{
			Name:        "build",
			Alias:       "b",
			Command:     packageManager + " run build",
			Description: "Build the project",
			Requires:    []string{"install"},
			Cache:       true,
			Inputs:      []string{"src/**/*", "package.json"},
			Outputs:     []string{"dist", "build"},
		}
	}

	cleanCommand := "rm -rf dist build node_modules .next out"
	if isWindows() {
		cleanCommand = "if exist dist rmdir /s /q dist & if exist build rmdir /s /q build & if exist node_modules rmdir /s /q node_modules & if exist .next rmdir /s /q .next & if exist out rmdir /s /q out"
	}

	cfg.Hooks["clean"] = models.Hook{
		Name:        "clean",
		Command:     cleanCommand,
		Description: "Clean build artifacts and dependencies",
	}

	return *cfg, nil
}

func readPackageJSON() (*packageJSON, error) {
	data, err := os.ReadFile("package.json")
	if err != nil {
		return nil, err
	}

	var pkg packageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	return &pkg, nil
}

func detectPackageManager() string {
	if _, err := os.Stat("pnpm-lock.yaml"); err == nil {
		return "pnpm"
	}
	if _, err := os.Stat("yarn.lock"); err == nil {
		return "yarn"
	}
	if _, err := os.Stat("bun.lockb"); err == nil {
		return "bun"
	}
	return "npm"
}

func getPackageLockFile(pm string) string {
	switch pm {
	case "pnpm":
		return "pnpm-lock.yaml"
	case "yarn":
		return "yarn.lock"
	case "bun":
		return "bun.lockb"
	default:
		return "package-lock.json"
	}
}

func determineDefaultTask(scripts map[string]string) string {
	if _, hasBuild := scripts["build"]; hasBuild {
		return "build"
	}
	if _, hasDev := scripts["dev"]; hasDev {
		return "dev"
	}
	if _, hasStart := scripts["start"]; hasStart {
		return "start"
	}
	return "build"
}

func generateTaskFromScript(name, command, pm string) models.Task {
	taskMap := map[string]struct {
		alias       string
		description string
		requires    []string
		inputs      []string
		outputs     []string
		cache       bool
		watch       bool
	}{
		"build": {
			alias:       "b",
			description: "Build the project",
			requires:    []string{"install"},
			inputs:      []string{"src/**/*", "package.json"},
			outputs:     []string{"dist", "build"},
			cache:       true,
		},
		"dev": {
			alias:       "d",
			description: "Start development server",
			requires:    []string{"install"},
			inputs:      []string{"src/**/*"},
			watch:       true,
		},
		"start": {
			alias:       "s",
			description: "Start the application",
			requires:    []string{"install"},
		},
		"test": {
			alias:       "t",
			description: "Run tests",
			requires:    []string{"install"},
			inputs:      []string{"src/**/*", "test/**/*", "**/*.test.*", "**/*.spec.*"},
			cache:       true,
		},
		"lint": {
			alias:       "l",
			description: "Lint code",
			requires:    []string{"install"},
			inputs:      []string{"src/**/*", ".eslintrc*", "eslint.config.*"},
			cache:       true,
		},
		"format": {
			alias:       "f",
			description: "Format code",
			inputs:      []string{"src/**/*", ".prettierrc*", "prettier.config.*"},
		},
		"type-check": {
			description: "Type check TypeScript",
			requires:    []string{"install"},
			inputs:      []string{"src/**/*.ts", "src/**/*.tsx", "tsconfig.json"},
			cache:       true,
		},
	}

	config, exists := taskMap[name]
	if !exists {
		return models.Task{
			Name:        name,
			Command:     pm + " run " + name,
			Description: "Run " + name + " script",
		}
	}

	return models.Task{
		Name:        name,
		Alias:       config.alias,
		Command:     pm + " run " + name,
		Description: config.description,
		Requires:    config.requires,
		Inputs:      config.inputs,
		Outputs:     config.outputs,
		Cache:       config.cache,
		Watch:       config.watch,
	}
}

func generateDefaultNode(cfg *config.Config) (config.Config, error) {
	cfg.DefaultTask = "build"

	cfg.Hooks["install"] = models.Hook{
		Name:        "install",
		Command:     "npm install",
		Description: "Install dependencies",
	}

	cfg.Tasks["build"] = models.Task{
		Name:        "build",
		Alias:       "b",
		Command:     "npm run build",
		Description: "Build the project",
		Requires:    []string{"install"},
		Cache:       true,
		Inputs:      []string{"src/**/*", "package.json"},
		Outputs:     []string{"dist"},
	}

	cfg.Tasks["test"] = models.Task{
		Name:        "test",
		Alias:       "t",
		Command:     "npm test",
		Description: "Run tests",
		Requires:    []string{"install"},
		Cache:       true,
		Inputs:      []string{"src/**/*", "test/**/*", "**/*.test.js", "**/*.spec.js"},
	}

	cfg.Tasks["dev"] = models.Task{
		Name:        "dev",
		Alias:       "d",
		Command:     "npm run dev",
		Description: "Start development server",
		Requires:    []string{"install"},
		Watch:       true,
		Inputs:      []string{"src/**/*"},
	}

	cleanCommand := "rm -rf dist node_modules"
	if isWindows() {
		cleanCommand = "if exist dist rmdir /s /q dist & if exist node_modules rmdir /s /q node_modules"
	}

	cfg.Hooks["clean"] = models.Hook{
		Name:        "clean",
		Command:     cleanCommand,
		Description: "Clean build artifacts",
	}

	return *cfg, nil
}

func isWindows() bool {
	return filepath.Separator == '\\'
}
