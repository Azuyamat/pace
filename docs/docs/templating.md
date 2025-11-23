---
sidebar_position: 6
---

# Project Templating

Pace can automatically generate project-specific configuration files tailored to your technology stack. The `pace init` command detects your project type and creates an optimized `config.pace` file with tasks, hooks, and sensible defaults.

## Quick Start

Navigate to your project directory and run:

```bash
pace init
```

Pace will:
1. Detect your project type (Go, Node.js, Python, Rust)
2. Analyze your project structure and configuration files
3. Generate a `config.pace` file with appropriate tasks
4. Update `.gitignore` to exclude `.pace-cache/`

## Manual Project Type Selection

If auto-detection doesn't work or you want to specify the project type manually:

```bash
pace init --type go
pace init --type node
pace init --type python
pace init --type rust
```

## Supported Project Types

### Node.js/TypeScript Projects

**Auto-detects:**
- Package manager (npm, pnpm, yarn, bun) from lock files
- All scripts in `package.json`
- TypeScript configuration

**Generated tasks include:**
- `install` - Install dependencies using detected package manager
- All scripts from `package.json` mapped as tasks
- Common scripts (build, dev, test, lint, format) with smart defaults
- Type-check task for TypeScript projects

**Example output:**

```pace
default build

task install [i] {
    command "pnpm install"
    description "Install dependencies"
    inputs ["package.json", "pnpm-lock.yaml"]
    outputs ["node_modules"]
    cache true
}

task build [b] {
    command "pnpm run build"
    description "Build the project"
    requires ["install"]
    inputs ["src/**/*", "package.json"]
    outputs ["dist", "build"]
    cache true
}

task dev [d] {
    command "pnpm run dev"
    description "Start development server"
    requires ["install"]
    inputs ["src/**/*"]
    watch true
}

task test [t] {
    command "pnpm run test"
    description "Run tests"
    requires ["install"]
    inputs ["src/**/*", "test/**/*", "**/*.test.*", "**/*.spec.*"]
    cache true
}
```

### Python Projects

**Auto-detects:**
- Dependency manager (pip, poetry, pdm, pipenv)
- Test framework (pytest, unittest)
- Linter (ruff, flake8, pylint)
- Formatter (ruff, black)
- Entry point (main.py, app.py, etc.)

**Generated tasks include:**
- `install` - Install dependencies using detected tool
- `run` - Run application (if entry point detected)
- `test` - Run tests with detected framework
- `lint` - Lint code with detected linter
- `format` - Format code with detected formatter
- `type-check` - Type checking with mypy (if pyproject.toml exists)

**Example output:**

```pace
default run

task install [i] {
    command "poetry install"
    description "Install dependencies"
    inputs ["pyproject.toml", "poetry.lock"]
    cache true
}

task run [r] {
    command "python main.py"
    description "Run the application"
    requires ["install"]
    inputs ["**/*.py"]
}

task test [t] {
    command "pytest"
    description "Run tests"
    requires ["install"]
    inputs ["src/**/*.py", "tests/**/*.py", "test/**/*.py"]
    cache true
}

task lint [l] {
    command "ruff check ."
    description "Lint code"
    requires ["install"]
    inputs ["**/*.py", "ruff.toml", ".ruff.toml"]
    cache true
}

task format [f] {
    command "ruff format ."
    description "Format code"
    inputs ["**/*.py", "ruff.toml", ".ruff.toml"]
}

task type-check {
    command "mypy ."
    description "Type check with mypy"
    requires ["install"]
    inputs ["**/*.py", "pyproject.toml"]
    cache true
}
```

### Rust Projects

**Auto-detects:**
- Project type (binary, library, or both)
- Package name from Cargo.toml
- Workspace configuration and members
- Benchmarks directory
- Examples directory

**Generated tasks include:**
- `build` - Build the project
- `release` - Build release binary
- `test` - Run tests
- `check` - Fast syntax checking
- `lint` - Clippy with strict warnings
- `format` - Format code with rustfmt
- `format-check` - Check formatting (for CI)
- `run` - Run binary (only for binary projects)
- `bench` - Run benchmarks (if benches/ exists)
- `examples` - Build examples (if examples/ exists)

**Example output:**

```pace
default build

task build [b] {
    command "cargo build"
    description "Build the project"
    inputs ["src/**/*.rs", "Cargo.toml", "Cargo.lock"]
    outputs ["target/debug"]
    cache true
}

task release [r] {
    command "cargo build --release"
    description "Build release binary"
    requires ["test"]
    inputs ["src/**/*.rs", "Cargo.toml", "Cargo.lock"]
    outputs ["target/release"]
    cache true
}

task test [t] {
    command "cargo test"
    description "Run tests"
    inputs ["src/**/*.rs", "tests/**/*.rs"]
    cache true
}

task lint [l] {
    command "cargo clippy -- -D warnings"
    description "Lint code with clippy"
    inputs ["src/**/*.rs"]
    cache true
}

task run {
    command "cargo run"
    description "Run the application"
    depends-on ["build"]
    inputs ["src/**/*.rs"]
    watch true
}
```

### Go Projects

**Auto-detects:**
- Module name from go.mod
- Project structure (standard, cmd-based, simple, library)
- Main package location (handles cmd/* patterns)
- golangci-lint configuration
- tools.go pattern for dev tools
- go:generate directives

**Generated tasks include:**
- `build` - Build the application with correct output path
- `test` - Run tests
- `lint` - Run golangci-lint (if config exists)
- `run` - Run the application
- `tools` - Install dev tools (if tools/tools.go exists)
- `generate` - Run go generate (if directives found)
- `tidy` - Tidy go.mod and go.sum
- `vet` - Run go vet

**Example output:**

```pace
default build

task build [b] {
    command "go build -o bin/myapp ./cmd/server"
    description "Build the application"
    requires ["test"]
    inputs ["**/*.go", "go.mod", "go.sum"]
    outputs ["bin/myapp"]
    cache true
}

task test [t] {
    command "go test ./..."
    description "Run tests"
    inputs ["**/*.go"]
    cache true
}

task lint [l] {
    command "golangci-lint run"
    description "Lint code"
    inputs ["**/*.go", ".golangci.yml", ".golangci.yaml"]
    cache true
}

task run [r] {
    command "go run ./cmd/server"
    description "Run the application"
    inputs ["**/*.go"]
    watch true
}

task generate {
    command "go generate ./..."
    description "Run go generate"
    inputs ["**/*.go"]
}

task tidy {
    command "go mod tidy"
    description "Tidy go.mod and go.sum"
    inputs ["go.mod"]
}

task vet {
    command "go vet ./..."
    description "Run go vet"
    inputs ["**/*.go"]
    cache true
}
```

## Detection Logic

### Node.js/TypeScript
Detected by presence of:
- `package.json`
- `node_modules/` directory

Package manager detected from:
- `pnpm-lock.yaml` → pnpm
- `yarn.lock` → yarn
- `bun.lockb` → bun
- `package-lock.json` → npm (default)

### Python
Detected by presence of:
- `*.py` files
- `main.py`
- `src/` or `tests/` directories

Dependency manager detected from:
- `pyproject.toml` with `[tool.poetry]` → poetry
- `pyproject.toml` with `[tool.pdm]` → pdm
- `Pipfile` → pipenv
- `requirements.txt` → pip (default)

### Rust
Detected by presence of:
- `Cargo.toml`
- `src/` directory

Project type detected from:
- `src/main.rs` or `[[bin]]` → binary
- `src/lib.rs` or `[lib]` → library

### Go
Detected by presence of:
- `go.mod` or `go.sum`
- `main.go`
- `cmd/` or `pkg/` directories

Project structure detected from:
- `cmd/` + (`pkg/` or `internal/`) → standard layout
- `cmd/` only → cmd-based layout
- `main.go` at root → simple layout
- No main package → library

## Customizing Generated Configuration

The generated `config.pace` file is a starting point. You can:

1. **Add custom tasks** - Extend with project-specific workflows
2. **Modify commands** - Adjust to your project's needs
3. **Update file patterns** - Refine inputs/outputs for better caching
4. **Add environment variables** - Configure task execution context
5. **Set up hooks** - Add pre/post execution logic

Example customization:

```pace
# Generated config
task build [b] {
    command "go build -o bin/app ./..."
    cache true
}

# Customized with version info
var version = "1.0.0"

task build [b] {
    command "go build -ldflags '-X main.Version=${version}' -o bin/app ./..."
    inputs ["**/*.go", "go.mod"]
    outputs ["bin/app"]
    cache true
    env {
        CGO_ENABLED = 0
        GOOS = linux
    }
}
```

## Platform-Specific Commands

Generated clean hooks automatically adapt to your platform:

**Unix/Linux/macOS:**
```pace
hook clean {
    command "rm -rf bin dist node_modules"
}
```

**Windows:**
```pace
hook clean {
    command "if exist bin rmdir /s /q bin"
}
```

## Best Practices

1. **Review generated config** - Always check the generated file before committing
2. **Adjust file patterns** - Refine inputs/outputs based on your project structure
3. **Add descriptions** - Document custom tasks you add
4. **Version control** - Commit `config.pace` to your repository
5. **Update regularly** - Re-run `pace init` when your project structure changes significantly

## Troubleshooting

**Q: Pace doesn't detect my project type**
- Ensure required files are present (package.json, Cargo.toml, go.mod, etc.)
- Use `--type` flag to specify manually
- Check that you're in the project root directory

**Q: Wrong package manager detected**
- Ensure lock file for your preferred manager exists
- Node.js detection priority: pnpm → yarn → bun → npm

**Q: Generated tasks don't match my workflow**
- Edit `config.pace` to customize tasks
- Tasks are templates - modify to fit your needs
- Consider contributing improvements to detection logic

**Q: Some tasks are missing**
- Pace only generates tasks for detected features
- Add custom tasks manually as needed
- Some advanced features require manual configuration

## Examples

See the [Examples](examples.md) page for complete project configurations and common patterns.
