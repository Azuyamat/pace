# pace run

Execute a task defined in your `config.pace` file.

## Usage

```bash
pace run [task-name] [flags] [arguments]
```

## Arguments

- `task-name` - Name of the task to run (optional if default task is set)

## Examples

### Run the default task

```bash
pace run
```

This runs the task specified with `default task_name` in your config.

### Run a specific task

```bash
pace run build
```

### Run a task with arguments

```bash
pace run greet --name=World
```

For a task defined as:
```pace
task greet {
    args {
        required ["name"]
    }
    command "echo 'Hello, $name!'"
}
```

### Run with positional arguments

```bash
pace run echo hello world
```

For a task defined as:
```pace
task echo {
    command "echo $1 $2"
}
```

### Run a task using an alias

```pace
# In config.pace
alias b build
```

```bash
pace run b  # Same as: pace run build
```

## Task Execution Flow

When you run a task, Pace executes in this order:

1. **Before hooks/tasks** - Tasks or hooks listed in `before`
2. **Main task** - The task itself
3. **After hooks/tasks** - Tasks or hooks listed in `after`
4. **Success/Failure hooks** - Based on task result

Example:
```pace
task deploy {
    before ["test", "build"]
    command "./scripts/deploy.sh"
    after ["cleanup"]
    on_success ["notify_success"]
    on_failure ["rollback"]
}
```

Execution order:
1. test
2. build
3. deploy command
4. cleanup
5. notify_success (if succeeded) OR rollback (if failed)

## Caching

If a task has `cache true`:

```pace
task build {
    cache true
    inputs ["**/*.go"]
    outputs ["bin/app"]
}
```

Pace will:
1. Check if inputs have changed since last run
2. Check if outputs exist
3. Skip execution if nothing changed
4. Display "Task up to date" message

To force execution even with cache:

```bash
pace run build --force
```

## Dependencies

Tasks automatically execute their dependencies:

```pace
task deploy {
    before ["build"]
}

task build {
    before ["test"]
}
```

Running `pace run deploy` will execute: test → build → deploy

## Parallel Execution

Enable parallel execution of dependencies:

```pace
task all {
    before ["backend", "frontend"]
    parallel true
}
```

Running `pace run all` will run backend and frontend simultaneously.

## Environment Variables

Tasks can define environment variables:

```pace
task build {
    env {
        "CGO_ENABLED" "0"
        "GOOS" "linux"
    }
    command "go build"
}
```

These variables are available to the command and override system environment variables.

## Working Directory

Change the working directory for a task:

```pace
task frontend {
    working_dir "frontend"
    command "npm run build"
}
```

## Timeout and Retries

Configure timeouts and automatic retries:

```pace
task deploy {
    command "./scripts/deploy.sh"
    timeout "10m"
    retry 3
    retry_delay "5s"
}
```

If the task exceeds 10 minutes, it will be terminated. If it fails, it will retry up to 3 times with a 5-second delay between attempts.

## Silent Mode

Suppress command output:

```pace
task quiet {
    silent true
    command "echo 'This will not be displayed'"
}
```

## Continue on Error

Allow execution to continue even if a task fails:

```pace
task optional {
    continue_on_error true
    command "might_fail.sh"
}
```

## Watch Mode

If a task has `watch true`, it automatically runs in watch mode:

```pace
task dev {
    watch true
    inputs ["**/*.go"]
    command "go run main.go"
}
```

Running `pace run dev` will watch the inputs and re-run on changes.

Alternatively, use the `pace watch` command:
```bash
pace watch build
```

## Notes

- If no task name is provided, the default task runs (if configured)
- Task names are case-sensitive
- Arguments can be positional or named (using `--name=value`)
- Dependencies are executed in order, and only once per run
- Circular dependencies are detected and reported as errors

## See Also

- [pace list](list.md) - List all available tasks
- [pace watch](watch.md) - Watch files and re-run tasks on changes
- [Configuration Reference](../configuration.md) - Learn about all task options
