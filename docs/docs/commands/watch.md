# pace watch

Watch task input files and automatically re-run the task when changes are detected.

## Usage

```bash
pace watch [task-name]
```

## Arguments

- `task-name` - Name of the task to watch (optional if default task is set)

## How It Works

The `watch` command monitors all files matching the task's `inputs` patterns. When any of these files change, the task is automatically re-executed.

## Examples

### Watch the default task

```bash
pace watch
```

### Watch a specific task

```bash
pace watch build
```

For a task defined as:
```pace
task build {
    command "go build -o bin/app main.go"
    inputs ["**/*.go", "go.mod"]
}
```

This will watch all `.go` files and `go.mod`, re-running the build whenever they change.

### Development server

```pace
task dev {
    description "Run development server"
    command "go run main.go"
    inputs ["**/*.go"]
}
```

```bash
pace watch dev
```

This creates a live-reload development environment that restarts your server on code changes.

## Task Configuration

For watch to work effectively, define the `inputs` property in your task:

```pace
task build {
    command "go build -o bin/app main.go"
    inputs [
        "**/*.go",      # All Go files
        "go.mod",       # Go module file
        "go.sum"        # Go dependencies
    ]
}
```

### Watch Property

Alternatively, set `watch true` on a task to make it watch by default:

```pace
task dev {
    watch true
    inputs ["**/*.go"]
    command "go run main.go"
}
```

Now running `pace run dev` will automatically start watch mode.

## File Patterns

Watch supports glob patterns for matching files:

```pace
task frontend {
    inputs [
        "src/**/*.ts",        # All TypeScript files in src/
        "src/**/*.tsx",       # All TSX files
        "public/**/*",        # All files in public/
        "package.json",       # Package file
        "tsconfig.json"       # TypeScript config
    ]
    command "npm run build"
}
```

### Common Patterns

- `**/*.go` - All Go files recursively
- `src/**/*.ts` - All TypeScript files in src/ and subdirectories
- `*.json` - All JSON files in current directory
- `**/*` - All files recursively (use with caution)

## Behavior

### Initial Run

When you start `pace watch`, it immediately runs the task once, then watches for changes.

### Debouncing

File changes are debounced to prevent multiple rapid executions. If multiple files change in quick succession, the task runs only once after the changes stabilize.

### Dependencies

Watch monitors only the specified task's inputs. If the task has dependencies:

```pace
task build {
    before ["test"]
    inputs ["**/*.go"]
}
```

Running `pace watch build` will:
1. Run test (once initially)
2. Run build
3. Watch build's inputs
4. On change: run test → run build

### Cache Interaction

Watch mode respects task caching settings. If a task has `cache true`, it may skip execution if inputs haven't actually changed (e.g., if the file was saved without modifications).

## Use Cases

### Development Server

```pace
task dev {
    watch true
    inputs ["**/*.go"]
    command "go run cmd/server/main.go"
}
```

### Frontend Development

```pace
task frontend-dev {
    watch true
    working_dir "frontend"
    inputs ["src/**/*.{ts,tsx}", "public/**/*"]
    command "npm run dev"
}
```

### Documentation Building

```pace
task docs {
    watch true
    inputs ["docs/**/*.md"]
    command "mkdocs build"
}
```

### Test-Driven Development

```pace
task test-watch {
    watch true
    inputs ["**/*.go", "**/*_test.go"]
    command "go test -v ./..."
}
```

## Stopping Watch

Press `Ctrl+C` to stop the watch process.

## Notes

- Watch requires the task to have `inputs` defined
- If no inputs are specified, watch will report an error
- Watch mode runs indefinitely until manually stopped
- Only one task can be watched at a time per command
- File system events may vary slightly between operating systems
- Very large numbers of files may impact watch performance

## See Also

- [pace run](run.md) - Run tasks without watching
- [pace list](list.md) - List all available tasks
- [Configuration Reference](../configuration.md) - Learn about task inputs and outputs
