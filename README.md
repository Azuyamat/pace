# Pace

[![Go Version](https://img.shields.io/github/go-mod/go-version/azuyamat/pace)](https://go.dev/)
[![License](https://img.shields.io/github/license/azuyamat/pace)](LICENSE)
[![Release](https://img.shields.io/github/v/release/azuyamat/pace)](https://github.com/azuyamat/pace/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/azuyamat/pace)](https://goreportcard.com/report/github.com/azuyamat/pace)

A cross-platform task runner and build orchestration tool. Define tasks once, run them anywhere.

[![Pace Logo](./public/BannerPace-640x320.png)](./public/BannerPace-640x320.png)

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

### Windows

#### PowerShell (recommended)
```powershell
iwr -useb https://raw.githubusercontent.com/Azuyamat/pace/refs/heads/master/install.ps1 | iex
```

#### Winget
```bash
winget install Azuyamat.Pace
```

### Linux

#### Script (recommended)
```bash
curl -sSL https://raw.githubusercontent.com/Azuyamat/pace/refs/heads/master/install.sh | sh
```

#### .deb Package
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

### Auto-generate Configuration (Recommended)

Let Pace automatically generate a configuration for your project:

```bash
cd your-project
pace init
```

Pace will detect your project type (Go, Node.js, Python, Rust) and create an optimized `config.pace` with appropriate tasks, caching, and file patterns.

### Manual Configuration

Alternatively, create a `config.pace` file in your project root:

```pace
default build

task build [b] {
    description "Build the project"
    command "go build -o bin/app main.go"
    requires [test]
    inputs [**/*.go]
    outputs [bin/app]
    cache true
}

hook test {
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
var output = "bin/myapp"
var version = "1.0.0"

task build {
    command "go build -ldflags '-X main.Version=${version}' -o ${output} ."
}
```

### Aliases

Create shortcuts for tasks using inline syntax:

```pace
task build [b] {
    description "Build the application"
    command "go build -o bin/app main.go"
}

task test [t] {
    description "Run tests"
    command "go test ./..."
}

task dev [d] {
    description "Start development server"
    command "go run main.go"
}
```

Or use standalone alias statements:

```pace
alias b build
alias t test
alias d dev
```

Now you can run:

```bash
pace run b    # same as: pace run build
```

### Task Properties

```pace
task example [ex] {
    description "Example task with all properties"
    command "echo 'Hello World'"

    # Dependencies and hooks
    depends-on [other-task]    # Tasks that must complete before this one
    requires [setup]           # Hooks to run before task
    triggers [cleanup]         # Hooks to run after task
    on_success [notify]        # Hooks to run on success
    on_failure [alert]         # Hooks to run on failure

    # File tracking for caching
    inputs [src/**/*.go, go.mod]
    outputs [bin/app]

    # Performance and behavior
    cache true                 # Enable smart caching
    timeout "10m"              # Maximum execution time
    retry 3                    # Number of retry attempts
    retry_delay "5s"           # Delay between retries

    # Execution environment
    working_dir "src"          # Working directory for command
    env {
        NODE_ENV = production
        DEBUG = false
    }

    # Conditional execution
    when "platform == linux"   # Only run on Linux

    # Flags
    silent false               # Suppress output
    parallel false             # Run dependencies in parallel
    continue_on_error false    # Continue if task fails
    watch false                # Enable file watching
}
```

**Property Defaults:**
- `cache`: false
- `parallel`: false
- `silent`: false
- `continue_on_error`: false
- `watch`: false
- `retry`: 0

**Note:** You can use identifiers or strings in arrays. Arrays support both `[item1, item2]` and `["item1", "item2"]` syntax.

### Arguments

Pass arguments to tasks:

```pace
task greet {
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
task echo {
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
default all

task all {
    description "Build all components"
    depends-on [backend, frontend]
    command "echo 'Build complete'"
}

task backend [be] {
    description "Build Go backend"
    command "go build -o bin/server cmd/server/main.go"
    inputs [cmd/**/*.go, internal/**/*.go]
    outputs [bin/server]
    cache true
}

task frontend [fe] {
    description "Build React frontend"
    command "npm run build"
    working_dir "frontend"
    inputs [frontend/src/**/*]
    outputs [frontend/dist/**/*]
    cache true
}
```

### Development Workflow

```pace
var app_name = "myapp"

task dev [d] {
    description "Start development server with hot reload"
    requires [generate]
    command "go run cmd/server/main.go"
    watch true
    inputs [**/*.go]
}

hook generate {
    description "Generate code"
    command "go generate ./..."
}

task test [t] {
    description "Run all tests"
    command "go test -v ./..."
    inputs [**/*.go]
    cache true
}

task lint [l] {
    description "Lint Go code"
    command "golangci-lint run"
    inputs [**/*.go]
    cache true
}

task build [b] {
    description "Build production binary"
    command "go build -o bin/${app_name} cmd/server/main.go"
    requires [test, lint]
    inputs [**/*.go]
    outputs [bin/${app_name}]
    cache true
}
```

## Contributing

Contributions are welcome! This is a passion project, so new ideas and improvements are always appreciated.

## License

MIT
