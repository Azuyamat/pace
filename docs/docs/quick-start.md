# Quick Start Guide

This guide will help you get started with Pace in just a few minutes.

## Your First Task

Let's add your first task:

```bash
pace add "Learn how to use Pace"
```

You should see a confirmation that your task was added successfully.

## Viewing Your Tasks

To see all your tasks:

```bash
pace list
```

For a tree view (default):

```bash
pace list --tree
```

## Completing Tasks

When you finish a task, mark it as complete:

```bash
pace complete 1
```

Replace `1` with the ID of the task you want to complete.

## Deleting Tasks

To remove a task entirely:

```bash
pace delete 1
```

## Common Workflows

### Daily Task Management

```bash
# Start your day by listing tasks
pace list

# Add new tasks as they come up
pace add "Review pull requests"
pace add "Update documentation"

# Complete tasks as you finish them
pace complete 1
pace complete 2

# End of day - check what's left
pace list
```

### Project-Based Tasks

```bash
# Add project tasks with descriptive names
pace add "Setup database schema"
pace add "Implement user authentication"
pace add "Write unit tests"

# View in tree format
pace list --tree
```

## Next Steps

- Learn about all available [commands](commands/add.md)
- Explore advanced features in the command reference
