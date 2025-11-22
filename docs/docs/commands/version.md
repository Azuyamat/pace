# version

Display version information for pace.

## Usage

```bash
pace version
```

## Description

The `version` command displays detailed version information about your pace installation, including:

- Version number
- Git commit hash
- Build date

This is useful for:
- Verifying your installation
- Reporting bugs
- Checking if you need to update

## Example

```bash
$ pace version
pace version 1.0.0
commit: a1b2c3d
built at: 2025-01-15T10:30:00Z
```

## Output

The command displays three lines of information:

1. **Version**: The semantic version number (e.g., 1.0.0)
2. **Commit**: The Git commit hash from which this version was built
3. **Built at**: The timestamp when the binary was compiled

## Development Builds

When running a development build (not installed via release), the version will show as:

```bash
pace version dev
commit: none
built at: unknown
```

## See Also

- [update](./update.md) - Update pace to the latest version
- [Installation](../installation.md) - How to install pace
