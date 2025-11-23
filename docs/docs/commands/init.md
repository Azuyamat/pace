---
sidebar_position: 1
---

# pace init

Initialize a new Pace project by automatically generating a `config.pace` file tailored to your project type.

## Usage

```bash
pace init [flags]
```

## Description

The `init` command analyzes your project and generates an optimized `config.pace` file with:
- Detected project type (Go, Node.js, Python, Rust)
- Appropriate tasks for your technology stack
- Smart defaults for caching, file patterns, and dependencies
- Hooks for common operations
- Updated `.gitignore` to exclude `.pace-cache/`

## Flags

### `--type, -t`

Manually specify the project type instead of auto-detection.

**Values:** `go`, `node`, `python`, `rust`

**Example:**
```bash
pace init --type go
pace init -t python
```

## Examples

### Auto-detect Project Type

```bash
cd my-project
pace init
```

**Output:**
```
INFO  Detected project type: go
INFO  Module: github.com/user/myapp
INFO  Project structure: standard
INFO  Main package: cmd/server
INFO  Generating template for go project...
INFO  --------------------------------------------------
INFO
default build

task build [b] {
    command "go build -o bin/myapp ./cmd/server"
    description "Build the application"
    ...
}
...
INFO  --------------------------------------------------
INFO  Generated template for go project
INFO  You may have to modify the generated config file according to your project structure.
INFO  Generated config.pace
INFO  Updated .gitignore to exclude .pace-cache/
```

### Specify Project Type

```bash
pace init --type node
```

Useful when:
- Auto-detection fails
- Project is not in standard structure
- You want to generate config for a specific language

### Re-initialize Project

Running `pace init` in a directory with an existing `config.pace` will **overwrite** it:

```bash
pace init  # Overwrites existing config.pace
```

**Best Practice:** Commit your current `config.pace` before re-running init, so you can review changes.

## What Gets Detected

### Node.js Projects
- Package manager (npm, pnpm, yarn, bun)
- All scripts from `package.json`
- TypeScript configuration
- Common frameworks (Next.js, React, Vue, etc.)

### Python Projects
- Dependency manager (pip, poetry, pdm, pipenv)
- Test framework (pytest, unittest)
- Linter (ruff, flake8, pylint)
- Formatter (ruff, black)
- Entry point files

### Rust Projects
- Project type (binary, library, both)
- Package name
- Workspace configuration
- Benchmarks and examples

### Go Projects
- Module name
- Project structure (standard, cmd-based, simple, library)
- Main package location
- Development tools (golangci-lint, go:generate)

## Generated Files

### `config.pace`

The main configuration file with all detected tasks and hooks.

### `.gitignore` (updated)

The `.pace-cache/` entry is automatically added if `.gitignore` exists.

## Post-Generation Steps

1. **Review the generated config** - Check that tasks match your workflow
2. **Customize as needed** - Add project-specific tasks or modify existing ones
3. **Test the tasks** - Run `pace list` and try executing tasks
4. **Commit to version control** - Add `config.pace` to your repository

## Common Issues

### "unsupported project type: unknown"

**Cause:** Pace couldn't detect your project type.

**Solutions:**
- Ensure required files exist (package.json, Cargo.toml, go.mod, etc.)
- Use `--type` flag to specify manually
- Verify you're in the project root directory

### Wrong package manager detected

**Cause:** Detection based on lock files.

**Solutions:**
- Delete unwanted lock files
- Node.js priority: pnpm → yarn → bun → npm
- Generate lock file for your preferred manager

### Generated config doesn't match my structure

**Cause:** Non-standard project layout.

**Solutions:**
- Edit `config.pace` manually to adjust paths
- Modify task commands to match your structure
- Update `inputs` and `outputs` patterns

### Tasks are missing features I use

**Cause:** Pace generates common tasks only.

**Solutions:**
- Add custom tasks for advanced features
- Modify generated tasks to include your flags
- Contribute detection improvements to Pace

## See Also

- [Project Templating](../templating.md) - Detailed templating documentation
- [Configuration Reference](../configuration.md) - Customize your config
- [Examples](../examples.md) - Sample configurations
