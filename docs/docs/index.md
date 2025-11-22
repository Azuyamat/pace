---
slug: /
sidebar_position: 1
---

# Welcome to Pace

Pace is a cross-platform task runner and build orchestration tool built with Go. Define your build tasks, development workflows, and automation scripts in a simple, human-readable configuration format that works the same on Windows, Linux, and macOS.

## Why Pace?

- **Cross-Platform**: Works identically on Windows, Linux, and macOS
- **Simple Syntax**: Human-readable `.pace` configuration files
- **Task Dependencies**: Chain tasks together with automatic ordering
- **File Watching**: Re-run tasks when source files change
- **Smart Caching**: Skip unchanged tasks based on file hashes
- **Hooks**: Run setup/cleanup commands before and after tasks
- **Parallel Execution**: Run independent tasks concurrently
- **Flexible**: Works with any language or build tool

If you've ever wanted `npm run` scripts but for Go, Rust, C++, or any other language—Pace is for you. Unlike Make (which has limited Windows support) or language-specific runners, Pace works the same everywhere.

## Quick Example

Create a `config.pace` file:

```pace
default build

task build {
    description "Build the application"
    command "go build -o bin/app main.go"
    inputs ["**/*.go"]
    outputs ["bin/app"]
    cache true
}

task test {
    description "Run tests"
    command "go test ./..."
    before ["build"]
}

task dev {
    description "Development server with auto-reload"
    watch true
    inputs ["**/*.go"]
    command "go run main.go"
}
```

Run your tasks:

```bash
# Run default task
pace run

# Run specific task
pace run test

# Watch and auto-reload
pace watch dev

# List all tasks
pace list
```

## Getting Started

New to Pace? Here's where to start:

1. **[Installation](installation.md)** - Install Pace on your system
2. **[Quick Start Guide](quick-start.md)** - Learn the basics in 5 minutes
3. **[Configuration Reference](configuration.md)** - Complete configuration guide
4. **[Commands Reference](commands/list.md)** - Explore all available commands
5. **[Examples](examples.md)** - Practical examples for different use cases

## Key Features

### Task Dependencies

Chain tasks together with automatic ordering:

```pace
task deploy {
    before ["test", "build"]
    command "./scripts/deploy.sh"
}
```

### Smart Caching

Skip tasks when nothing has changed:

```pace
task build {
    cache true
    inputs ["src/**/*.go"]
    outputs ["bin/app"]
}
```

### File Watching

Auto-reload during development:

```pace
task dev {
    watch true
    inputs ["**/*.go"]
    command "go run main.go"
}
```

### Multi-Language Support

Works with any language or tool:

```pace
task backend {
    command "go build -o bin/server main.go"
}

task frontend {
    working_dir "frontend"
    command "npm run build"
}
```

## Need Help?

- Check out our [GitHub repository](https://github.com/azuyamat/pace)
- Open an issue for bug reports or feature requests
- Contribute to the project - PRs are welcome!
- Star the project if you find it useful
