# Pace

[![Go Version](https://img.shields.io/github/go-mod/go-version/azuyamat/pace)](https://go.dev/)
[![License](https://img.shields.io/github/license/azuyamat/pace)](LICENSE)
[![Release](https://img.shields.io/github/v/release/azuyamat/pace)](https://github.com/azuyamat/pace/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/azuyamat/pace)](https://goreportcard.com/report/github.com/azuyamat/pace)

A cross-platform task runner and build orchestration tool. Define tasks once, run them anywhere.

## Why Pace?

If you've ever wanted `npm run` scripts but for Go, Rust, C++, or any other language—Pace is for you. Unlike Make (which has limited Windows support) or language-specific runners, Pace works the same on Windows, Linux, and macOS with an intuitive configuration syntax.

## Features

- **Cross-platform** - Works on Windows, Linux, and macOS
- **Simple syntax** - Human-readable `.pace` configuration files
- **Task dependencies** - Chain tasks together with automatic ordering
- **File watching** - Re-run tasks when source files change
- **Smart caching** - Skip unchanged tasks based on file hashes
- **Hooks** - Run setup/cleanup commands before and after tasks
- **Arguments** - Pass parameters to tasks with positional or named args
- **Retry logic** - Automatically retry failed tasks
- **Parallel execution** - Run independent tasks concurrently

## Installation

### Windows (winget)

```bash
winget install Azuyamat.Pace
```

### Linux (deb)

Download the `.deb` package from [releases](https://github.com/azuyamat/pace/releases):

```bash
sudo dpkg -i pace_<version>_amd64.deb
```

### Go

```bash
go install github.com/azuyamat/pace/cmd/pace@latest
```

### From Releases

Download the appropriate binary for your platform from the [releases page](https://github.com/azuyamat/pace/releases).

## Quick Start

Create a `config.pace` file in your project root:

```pace
default "build"

task "build" {
    description "Build the project"
    command "go build -o bin/app main.go"
    before ["test"]
    inputs ["**/*.go"]
    outputs ["bin/app"]
}

hook "test" {
    description "Run tests"
    command "go test ./..."
}
```

Run your default task:

```bash
pace run
```

Or run a specific task:

```bash
pace run build
```

## Configuration

### Variables

Define reusable values:

```pace
set "output" "bin/myapp"
set "version" "1.0.0"

task "build" {
    command "go build -ldflags '-X main.Version=${version}' -o ${output} ."
}
```

### Aliases

Create shortcuts for tasks:

```pace
alias "b" "build"
alias "t" "test"
alias "d" "dev"
```

Now you can run:

```bash
pace run b    # same as: pace run build
```

### Task Properties

```pace
task "example" {
    description "Example task"
    command "echo 'Hello World'"

    before ["setup"]
    after ["cleanup"]
    on_success ["notify"]
    on_failure ["alert"]

    inputs ["src/**/*.go", "go.mod"]
    outputs ["bin/app"]

    cache true
    timeout "10m"
    retry 3
    retry_delay "5s"

    working_dir "src"
    env {
        "NODE_ENV" "production"
        "DEBUG" "false"
    }

    silent false
    parallel false
    continue_on_error false
}
```

### Arguments

Pass arguments to tasks:

```pace
task "greet" {
    args {
        required ["name"]
    }
    command "echo 'Hello, $name!'"
}
```

```bash
pace run greet --name=World
```

Or use positional arguments:

```pace
task "echo" {
    command "echo $1 $2"
}
```

```bash
pace run echo hello world
```

### Hooks

Hooks are lightweight tasks for setup/cleanup:

```pace
hook "format" {
    description "Format code"
    command "gofmt -s -w ."
}
```

## Commands

```bash
pace run [task]      # Run a task (or default task)
pace watch [task]    # Watch inputs and re-run on changes
pace list            # List all tasks and hooks
pace list --tree     # List with dependency tree
pace help [command]  # Show help
```

### Flags

- `--dry-run` - Show what would run without executing
- `--force` - Ignore cache and force execution

## File Watching

Watch task inputs and automatically re-run:

```bash
pace watch build
```

This monitors all files matching the task's `inputs` patterns and re-executes when changes are detected.

## Caching

When `cache true` is set, Pace tracks:

- Input file hashes
- Output file hashes
- Command string
- Dependency results

If nothing has changed since the last run, the task is skipped. Cache data is stored in `.pace-cache/`.

## VS Code Extension

Syntax highlighting for `.pace` files is available. Check the [vscode-pace](vscode-pace/) directory for the extension.

## Examples

### Multi-language Project

```pace
default "all"

task "all" {
    before ["backend", "frontend"]
    command "echo 'Build complete'"
}

task "backend" {
    command "go build -o bin/server cmd/server/main.go"
    inputs ["cmd/**/*.go", "internal/**/*.go"]
    outputs ["bin/server"]
    cache true
}

task "frontend" {
    command "npm run build"
    working_dir "frontend"
    inputs ["frontend/src/**/*"]
    outputs ["frontend/dist/**/*"]
    cache true
}
```

### Development Workflow

```pace
task "dev" {
    description "Start development server"
    before ["generate"]
    command "go run cmd/server/main.go"
}

hook "generate" {
    command "go generate ./..."
}

task "test" {
    command "go test -v ./..."
    inputs ["**/*.go"]
    cache true
}

task "lint" {
    command "golangci-lint run"
    inputs ["**/*.go"]
}
```

## Contributing

Contributions are welcome! This is a passion project, so new ideas and improvements are always appreciated.

## License

MIT
