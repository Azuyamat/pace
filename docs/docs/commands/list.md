# pace list

Display all your tasks.

## Usage

```bash
pace list [flags]
```

## Flags

- `--tree` - Display tasks in a tree view format (default: true)

## Examples

### List all tasks (tree view)

```bash
pace list
```

### List all tasks with explicit tree flag

```bash
pace list --tree
```

### List without tree view

```bash
pace list --tree=false
```

## Output Format

### Tree View

When using tree view, tasks are displayed with:
- Task ID
- Completion status (✓ for completed)
- Task name

Example output:
```
├── [1] ✓ Task 1
├── [2] Task 2
└── [3] Task 3
```

## Notes

- Completed tasks are marked with a checkmark (✓)
- Tasks are displayed in the order they were added
- The tree view provides a clean, visual representation of your task list
