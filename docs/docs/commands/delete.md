# pace delete

Delete a task from your task list.

## Usage

```bash
pace delete <task-id>
```

## Arguments

- `task-id` - The ID of the task to delete (required)

## Examples

### Delete a single task

```bash
pace delete 1
```

### Delete multiple tasks

```bash
pace delete 1
pace delete 2
pace delete 3
```

## Finding Task IDs

To find the ID of a task, use the `list` command:

```bash
pace list
```

Task IDs are shown in square brackets next to each task.

## Notes

- This action is permanent and cannot be undone
- The task will be completely removed from your task list
- If you want to keep a record of completed tasks, use `pace complete` instead
- Deleting a non-existent task will show an error message
