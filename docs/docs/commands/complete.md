# pace complete

Mark a task as completed.

## Usage

```bash
pace complete <task-id>
```

## Arguments

- `task-id` - The ID of the task to mark as complete (required)

## Examples

### Complete a single task

```bash
pace complete 1
```

### Complete multiple tasks

```bash
pace complete 1
pace complete 2
pace complete 3
```

## Finding Task IDs

To find the ID of a task, use the `list` command:

```bash
pace list
```

Task IDs are shown in square brackets next to each task.

## Notes

- Completed tasks remain in your task list but are marked with a checkmark (✓)
- You can view completed tasks with `pace list`
- To remove a task entirely, use `pace delete` instead
- Completing an already completed task will show a success message
