# update

Update pace to the latest version.

## Usage

```bash
pace update
```

## Description

The `update` command automatically updates pace to the latest version available on GitHub. The command intelligently detects your installation method and uses the appropriate update mechanism.

## Installation Method Detection

Pace automatically detects how it was installed and updates accordingly:

### Winget (Windows)

If pace was installed via Winget, the command runs:

```bash
winget upgrade Azuyamat.Pace
```

### Go Install

If pace was installed via `go install`, the command runs:

```bash
go install github.com/azuyamat/pace/cmd/pace@latest
```

### Direct Installation

If pace was installed by downloading the binary directly (via curl, wget, or manual download), the command:

1. Downloads the latest release from GitHub
2. Extracts the binary from the archive
3. Safely replaces the current binary with a backup mechanism

## Examples

### Check and Update

```bash
$ pace update
Checking for updates...
Updating from 1.0.0 to 1.1.0...
Proceeding with self-update...

Successfully updated to version 1.1.0
Restart pace to use the new version
```

### Already Up to Date

```bash
$ pace update
Checking for updates...
Already on latest version: 1.1.0
```

### Winget Installation

```bash
$ pace update
Checking for updates...
Updating from 1.0.0 to 1.1.0...
Running: winget upgrade Azuyamat.Pace
[winget output]
Update completed successfully!
```

### Go Installation

```bash
$ pace update
Checking for updates...
Updating from 1.0.0 to 1.1.0...
Running: go install github.com/azuyamat/pace/cmd/pace@latest
[go install output]
Update completed successfully!
```

## How It Works

1. **Check for Updates**: Queries the GitHub API for the latest release
2. **Compare Versions**: Compares your current version with the latest
3. **Detect Installation Method**: Determines how pace was installed
4. **Execute Update**: Runs the appropriate update mechanism
5. **Verify**: Confirms the update was successful

## Safety Features

For direct installations, the update process includes:

- **Backup Creation**: The old binary is backed up before replacement
- **Rollback on Failure**: If the update fails, the backup is restored
- **Permission Checks**: Verifies write permissions before attempting update

## Platform Support

The update command works on:

- **Windows** (amd64, arm64)
- **Linux** (amd64, arm64)
- **macOS** (amd64, arm64)

## Troubleshooting

### Permission Denied

If you get a permission error:

- **Linux/macOS**: You may need to run with `sudo` if the binary is in a protected directory
- **Windows**: Run your terminal as Administrator if installed in Program Files

### Update Command Not Available

If Winget or Go is detected but the command fails:

- Ensure `winget` or `go` is installed and available in your PATH
- Verify you have an active internet connection

### Manual Update

If the automatic update fails, you can always:

1. Download the latest release from [GitHub Releases](https://github.com/azuyamat/pace/releases)
2. Extract the binary
3. Replace your current pace binary

## See Also

- [version](./version.md) - Check your current version
- [Installation](../installation.md) - Installation methods
