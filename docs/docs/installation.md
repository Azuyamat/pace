# Installation

## Prerequisites

- Go 1.16 or higher (for Go installation method)
- No runtime dependencies required

## Installation Methods

### Windows (winget)

The easiest way to install on Windows:

```bash
winget install Azuyamat.Pace
```

### Linux (deb)

Download the `.deb` package from [releases](https://github.com/azuyamat/pace/releases):

```bash
sudo dpkg -i pace_<version>_amd64.deb
```

### Go Install

Install the latest version using Go:

```bash
go install github.com/azuyamat/pace/cmd/pace@latest
```

This will install the `pace` binary to your `$GOPATH/bin` directory. Make sure `$GOPATH/bin` is in your system PATH.

### From Releases

Download pre-built binaries for your platform:

1. Visit the [releases page](https://github.com/azuyamat/pace/releases)
2. Download the appropriate binary for your platform
3. Extract and move to a directory in your PATH

**Platform-specific locations:**
- **Linux/macOS**: `/usr/local/bin/pace`
- **Windows**: `C:\Program Files\pace\pace.exe` (or any directory in your PATH)

### Building from Source

If you prefer to build from source:

```bash
# Clone the repository
git clone https://github.com/azuyamat/pace.git
cd pace

# Build the binary
go build -o pace ./cmd/pace/main.go

# Install to your PATH (Linux/macOS)
sudo mv pace /usr/local/bin/

# Or on Windows
move pace.exe C:\Program Files\pace\
```

## Verify Installation

After installation, verify that Pace is working correctly:

```bash
pace run
```

If you see an error about `config.pace` not found, that's normal - it means Pace is installed correctly and looking for a configuration file.

## Configuration File

Pace looks for a `config.pace` file in your current working directory. No global configuration is required. Each project can have its own `config.pace` file.

### Cache Directory

When caching is enabled, Pace stores cache data in `.pace-cache/` within your project directory. This directory should be added to your `.gitignore` file:

```gitignore
.pace-cache/
```

## Next Steps

- [Quick Start Guide](quick-start.md) - Create your first task configuration
- [Configuration Reference](configuration.md) - Learn about all available options
