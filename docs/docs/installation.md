# Installation

## Prerequisites

- Go 1.16 or higher

## Install from Source

The easiest way to install Pace is using `go install`:

```bash
go install github.com/azuyamat/pace@latest
```

This will install the `pace` binary to your `$GOPATH/bin` directory.

## Verify Installation

After installation, verify that Pace is working correctly:

```bash
pace --version
```

You should see the version information displayed.

## Building from Source

If you prefer to build from source:

```bash
# Clone the repository
git clone https://github.com/azuyamat/pace.git
cd pace

# Build the binary
go build -o pace ./cmd/pace

# Move to your PATH (optional)
sudo mv pace /usr/local/bin/
```

## Configuration

Pace stores its data in a local JSON file. By default, tasks are stored in:

- **Linux/macOS**: `~/.pace/tasks.json`
- **Windows**: `%USERPROFILE%\.pace\tasks.json`

No additional configuration is required to get started.
