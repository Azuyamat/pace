# pace list

Display all available tasks, hooks, and aliases from your `config.pace` file.

## Usage

```bash
pace list [flags]
```

## Flags

- `--tree`, `-t` - Display tasks in a tree view showing dependencies (default: false)

## Examples

### List all tasks

```bash
pace list
```

Example output:
```
Available tasks:

  build                Build the application (default)
  test                 Run all tests
  deploy               Deploy to production

Aliases:
  b                    -> build
  t                    -> test

Available hooks:
  format               Format code
  lint                 Run linter
```

### List with dependency tree

```bash
pace list --tree
```

Example output:
```
Task dependency tree:

build (default)
  ├── test
  └── lint

deploy
  └── build
      ├── test
      └── lint
```

### Short flag version

```bash
pace list -t
```

## Output Format

### Simple List (default)

The default view shows three sections:

1. **Available tasks**: All defined tasks with their descriptions
   - Tasks marked `(default)` will run when you execute `pace run` without arguments
   - If no description is provided, the command is shown instead

2. **Aliases**: Shortcuts to tasks (if any defined)

3. **Available hooks**: Reusable hooks that can be referenced by tasks

### Tree View

The tree view shows:
- Task names with their dependencies
- Visual tree structure showing the execution order
- Default task indicator
- Circular dependency detection (marked as "circular")

## Notes

- Tasks are displayed in alphabetical order
- The default task is marked with `(default)`
- Tree view helps visualize complex dependency chains
- Circular dependencies are detected and marked to prevent infinite loops
