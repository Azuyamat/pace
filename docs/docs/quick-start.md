# Quick Start Guide

This guide will help you get started with Pace in just a few minutes.

## Create Your First Configuration

### Option 1: Auto-generate (Recommended)

Let Pace automatically generate a configuration for your project:

```bash
pace init
```

This will detect your project type (Go, Node.js, Python, Rust) and create an optimized `config.pace` with appropriate tasks. See [Project Templating](templating.md) for details.

### Option 2: Manual Configuration

Create a file named `config.pace` in your project root:

```pace
# Set the default task
default build

# Define a build task with alias
task build [b] {
    description "Build the application"
    command "go build -o bin/app main.go"
    inputs [**/*.go]
    outputs [bin/app]
    cache true
}

# Define a test task
task test [t] {
    description "Run tests"
    command "go test ./..."
}
```

## Run Your First Task

Run the default task (build):

```bash
pace run
```

Or run a specific task:

```bash
pace run test
```

## List Available Tasks

See all available tasks and their descriptions:

```bash
pace list
```

For a tree view showing dependencies:

```bash
pace list --tree
```

## Add Dependencies

Tasks can depend on other tasks and hooks:

```pace
task build {
    description "Build the application"
    command "go build -o bin/app main.go"
    requires [test]  # Run test hook before build
}

hook test {
    description "Run tests"
    command "go test ./..."
}

task deploy {
    description "Deploy to production"
    command "./scripts/deploy.sh"
    depends-on [build]  # Ensure build completes before deploy
    requires [test]     # Run test before deploy
}
```

Now when you run `pace run deploy`, it will automatically run: test → build → deploy

## Use Variables

Define reusable values:

```pace
# Define variables
var output = "bin/myapp"
var version = "1.0.0"

task build {
    command "go build -ldflags '-X main.Version=${version}' -o ${output} main.go"
}
```

## Watch for Changes

Automatically re-run tasks when files change:

```bash
pace watch build
```

This will monitor all files matching the task's `inputs` patterns and re-execute when changes are detected.

## Add Task Aliases

Create shortcuts using inline syntax:

```pace
task build [b] {
    description "Build application"
    command "go build -o bin/app main.go"
}

task test [t] {
    description "Run tests"
    command "go test ./..."
}

task deploy [d] {
    description "Deploy to production"
    command "./scripts/deploy.sh"
}
```

Or use standalone alias statements:

```pace
alias b build
alias t test
alias d deploy
```

Now you can run:

```bash
pace run b    # same as: pace run build
pace run t    # same as: pace run test
```

## Common Workflows

### Development Workflow

```pace
default dev

task dev [d] {
    description "Start development server"
    command "go run main.go"
    requires [build]
}

task build [b] {
    description "Build the project"
    command "go build -o bin/app main.go"
    inputs [**/*.go]
    cache true
}

task test [t] {
    description "Run tests"
    command "go test -v ./..."
}
```

### Build and Test Pipeline

```pace
default all

task all {
    description "Run full pipeline"
    depends-on [lint, test, build]
}

task lint [l] {
    description "Run linter"
    command "golangci-lint run"
}

task test [t] {
    description "Run tests"
    command "go test ./..."
}

task build [b] {
    description "Build application"
    command "go build -o bin/app main.go"
    cache true
}
```

## Next Steps

- [Project Templating](templating.md) - Auto-generate configurations for your projects
- [Configuration Reference](configuration.md) - Learn about all configuration options
- [Commands Reference](commands/list.md) - Explore all available commands
- [Examples](examples.md) - See more practical examples
