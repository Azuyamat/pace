# Welcome to Pace

Pace is a cross-platform task runner and build orchestration tool built with Go. Define your build tasks, development workflows, and automation scripts in a simple, human-readable configuration format that works the same on Windows, Linux, and macOS.

## Why Pace?

If you've ever wanted `npm run` scripts but for Go, Rust, C++, or any other language—Pace is for you. Unlike Make (which has limited Windows support) or language-specific runners, Pace works the same on all platforms with an intuitive configuration syntax.

## Features

- **Cross-platform**: Works identically on Windows, Linux, and macOS
- **Simple syntax**: Human-readable `.pace` configuration files
- **Task dependencies**: Chain tasks together with automatic ordering
- **File watching**: Re-run tasks when source files change
- **Smart caching**: Skip unchanged tasks based on file hashes
- **Hooks**: Run setup/cleanup commands before and after tasks
- **Arguments**: Pass parameters to tasks with positional or named args
- **Retry logic**: Automatically retry failed tasks
- **Parallel execution**: Run independent tasks concurrently

## Quick Start

Create a `config.pace` file in your project root:

```pace
default build

task build {
    description "Build the project"
    command "go build -o bin/app main.go"
    inputs ["**/*.go"]
    outputs ["bin/app"]
    cache true
}

task test {
    description "Run tests"
    command "go test ./..."
}
```

Run your tasks:

```bash
# Run default task
pace run

# Run specific task
pace run build

# List all tasks
pace list
```

## Next Steps

- [Installation](installation.md) - Detailed installation instructions
- [Quick Start Guide](quick-start.md) - Learn the basics with examples
- [Configuration Reference](configuration.md) - Complete configuration guide
- [Commands Reference](commands/list.md) - All available commands
