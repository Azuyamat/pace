# Configuration Reference

This page provides a complete reference for Pace configuration files (`.pace` files).

## File Structure

Pace configuration files use a simple, human-readable syntax. The file typically named `config.pace` should be placed in your project root.

```pace
# Comments start with #

# Define variables
var variable_name = "value"

# Set default task
default task_name

# Create task aliases
alias shortname taskname

# Define tasks
task task_name {
    # Task properties...
}

# Define hooks
hook hook_name {
    # Hook properties...
}
```

## Variables

Variables allow you to define reusable values throughout your configuration.

### Syntax

```pace
var variable_name = "value"
```

### Variable Interpolation

Variables can be referenced using `${variable_name}` or `$variable_name` syntax:

```pace
var output = "bin/myapp"
var version = "1.0.0"

task build {
    command "go build -ldflags '-X main.Version=${version}' -o ${output} main.go"
}
```

### Environment Variables

Environment variables from your system are automatically available:

```pace
task build {
    command "echo Building on ${HOSTNAME}"
}
```

## Default Task

Set the default task that runs when you execute `pace run` without arguments:

```pace
default build
```

## Aliases

Create shortcuts for task names:

```pace
alias b build
alias t test
alias d deploy
```

Usage:
```bash
pace run b  # same as: pace run build
```

## Tasks

Tasks are the main building blocks of your Pace configuration.

### Basic Task

```pace
task build {
    description "Build the application"
    command "go build -o bin/app main.go"
}
```

### Task Properties

#### `description` (string)
Human-readable description of what the task does.

```pace
task build {
    description "Build the application for production"
}
```

#### `command` (string, required)
The command to execute. Can span multiple lines using triple quotes.

```pace
task build {
    command "go build -o bin/app main.go"
}

task multi {
    command """
        echo "Step 1"
        echo "Step 2"
        go build
    """
}
```

#### `inputs` (array of strings)
File patterns that this task depends on. Used for caching and file watching.

```pace
task build {
    inputs ["**/*.go", "go.mod", "go.sum"]
}
```

Supports glob patterns:
- `**/*.go` - All Go files recursively
- `src/**/*.ts` - All TypeScript files in src/
- `*.json` - All JSON files in current directory

#### `outputs` (array of strings)
Files or patterns that this task produces. Used for caching.

```pace
task build {
    outputs ["bin/app", "bin/app.exe"]
}
```

#### `cache` (boolean)
Enable smart caching. Task will be skipped if inputs haven't changed since last successful run.

```pace
task build {
    cache true
    inputs ["**/*.go"]
    outputs ["bin/app"]
}
```

Default: `false`

#### `dependencies` (array of strings)
Tasks that must complete before this task runs (alias for `before`).

```pace
task deploy {
    dependencies ["test", "build"]
}
```

#### `before` (array of strings)
Hooks or tasks to run before this task.

```pace
task build {
    before ["test", "lint"]
}
```

#### `after` (array of strings)
Hooks or tasks to run after this task completes.

```pace
task build {
    after ["notify"]
}
```

#### `on_success` (array of strings)
Hooks or tasks to run only if this task succeeds.

```pace
task deploy {
    on_success ["notify_success"]
}
```

#### `on_failure` (array of strings)
Hooks or tasks to run only if this task fails.

```pace
task deploy {
    on_failure ["rollback", "notify_failure"]
}
```

#### `env` (map of strings)
Environment variables to set for this task.

```pace
task build {
    env {
        "CGO_ENABLED" "0"
        "GOOS" "linux"
        "GOARCH" "amd64"
    }
}
```

#### `working_dir` (string)
Directory to run the command in.

```pace
task frontend {
    working_dir "frontend"
    command "npm run build"
}
```

#### `timeout` (string)
Maximum time to allow the task to run. Supports duration suffixes.

```pace
task build {
    timeout "10m"
}
```

Supported units:
- `s` - seconds
- `m` - minutes
- `h` - hours

Examples: `30s`, `5m`, `1h30m`

#### `retry` (integer)
Number of times to retry if the task fails.

```pace
task deploy {
    retry 3
}
```

Default: `0` (no retries)

#### `retry_delay` (string)
Time to wait between retries.

```pace
task deploy {
    retry 3
    retry_delay "5s"
}
```

#### `parallel` (boolean)
Whether dependencies can run in parallel.

```pace
task all {
    before ["backend", "frontend"]
    parallel true  # Run backend and frontend simultaneously
}
```

Default: `false`

#### `silent` (boolean)
Suppress command output.

```pace
task quiet {
    silent true
}
```

Default: `false`

#### `continue_on_error` (boolean)
Continue execution even if this task fails.

```pace
task optional {
    continue_on_error true
}
```

Default: `false`

#### `watch` (boolean)
Automatically watch input files and re-run on changes.

```pace
task dev {
    watch true
    inputs ["**/*.go"]
    command "go run main.go"
}
```

Default: `false`

#### `args` (object)
Define arguments that can be passed to the task.

```pace
task greet {
    args {
        required ["name"]
        optional ["greeting"]
    }
    command "echo '${greeting:-Hello}, $name!'"
}
```

Usage:
```bash
pace run greet --name=World --greeting=Hi
```

Positional arguments are also supported:
```pace
task echo {
    command "echo $1 $2 $3"
}
```

Usage:
```bash
pace run echo hello world test
```

## Hooks

Hooks are lightweight tasks designed for setup, cleanup, or other auxiliary operations.

### Basic Hook

```pace
hook format {
    description "Format code"
    command "gofmt -s -w ."
}
```

### Hook Properties

Hooks support a subset of task properties:

- `description` - Description of the hook
- `command` - Command to execute (required)
- `env` - Environment variables
- `working_dir` - Working directory

Hooks **do not** support:
- Dependencies
- Caching
- Retries
- Arguments
- Before/after hooks

## Globals

Define global settings that apply to all tasks.

```pace
globals {
    "env" {
        "GO_ENV" "production"
    }
}
```

## Imports

Import configuration from other files:

```pace
import "tasks/common.pace"
import "tasks/deploy.pace"
```

Imported configurations are merged with the current file. Local definitions take precedence.

## Complete Example

```pace
# Variables
var app_name = "myapp"
var version = "1.0.0"
var build_dir = "bin"

# Default task
default build

# Aliases
alias b build
alias t test
alias d deploy

# Hooks
hook format {
    description "Format Go code"
    command "gofmt -s -w ."
}

hook lint {
    description "Run linter"
    command "golangci-lint run"
}

# Tasks
task test {
    description "Run all tests"
    command "go test -v ./..."
    inputs ["**/*.go"]
    timeout "5m"
    cache true
}

task build {
    description "Build the application"
    command "go build -ldflags '-X main.Version=${version}' -o ${build_dir}/${app_name} main.go"
    before ["test", "format", "lint"]
    inputs ["**/*.go", "go.mod", "go.sum"]
    outputs ["${build_dir}/${app_name}"]
    cache true
    env {
        "CGO_ENABLED" "0"
    }
}

task docker {
    description "Build Docker image"
    command "docker build -t ${app_name}:${version} ."
    before ["build"]
    inputs ["Dockerfile", "${build_dir}/${app_name}"]
}

task deploy {
    description "Deploy to production"
    command "./scripts/deploy.sh ${version}"
    before ["docker"]
    on_success ["notify_success"]
    on_failure ["notify_failure", "rollback"]
    retry 2
    retry_delay "10s"
    timeout "15m"
}

hook notify_success {
    command "echo 'Deployment successful!'"
}

hook notify_failure {
    command "echo 'Deployment failed!'"
}

hook rollback {
    command "./scripts/rollback.sh"
}
```

## Comments

Comments start with `#` and continue to the end of the line:

```pace
# This is a comment
task build {  # Inline comments are also supported
    command "go build"
}
```

## String Literals

Pace supports both single-line and multi-line strings:

```pace
# Single-line
task single {
    command "echo 'hello'"
}

# Multi-line (using triple quotes)
task multi {
    command """
        echo "Line 1"
        echo "Line 2"
        echo "Line 3"
    """
}
```

## Boolean Values

Use `true` or `false` (lowercase):

```pace
task build {
    cache true
    silent false
    parallel true
}
```

## Numbers

Numbers can be integers or include duration suffixes:

```pace
task build {
    retry 3
    timeout "5m"
    retry_delay "30s"
}
```

## Next Steps

- [Commands Reference](commands/list.md) - Learn about all CLI commands
- [Examples](examples.md) - See practical configuration examples
