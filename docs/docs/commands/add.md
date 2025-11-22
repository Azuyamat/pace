# pace add

Add a new task to your task list.

## Usage

```bash
pace add <task-name>
```

## Arguments

- `task-name` - The name/description of the task to add (required)

## Examples

### Add a simple task

```bash
pace add "Write documentation"
```

### Add a task with multiple words

```bash
pace add "Review and merge pull request #42"
```

### Add multiple tasks

```bash
pace add "Task 1"
pace add "Task 2"
pace add "Task 3"
```

## Notes

- Task names can contain spaces and special characters
- Each task is automatically assigned a unique ID
- Tasks are stored in your local tasks.json file
- There is no limit to the number of tasks you can add
