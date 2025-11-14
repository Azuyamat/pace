# Pace Language Support for VS Code

Syntax highlighting support for Pace configuration files (`.pace`).

## Features

- **Syntax Highlighting**: Full syntax highlighting for Pace configuration files
- **Comment Support**: Line comments with `#`
- **Auto-closing**: Automatic closing of brackets, braces, and quotes
- **Code Folding**: Fold task, hook, and block definitions
- **Variable Highlighting**: Highlights variable interpolation (`${VAR}` and `$VAR`)

## Supported Syntax

- **Keywords**: `set`, `default`, `alias`, `globals`, `hook`, `task`
- **Properties**: `description`, `command`, `before`, `after`, `inputs`, `outputs`, `cache`, `on_success`, `on_failure`, `dependencies`, `parallel`, `timeout`, `retry`, `retry_delay`, `silent`, `watch`, `continue_on_error`, `env`, `args`, `required`, `optional`
- **Strings**: Single-line (`"..."`) and multi-line (`"""..."""`)
- **Variables**: Variable interpolation in strings
- **Comments**: Line comments with `#`
- **Booleans**: `true`, `false`
- **Numbers**: Including duration suffixes (`5m`, `30s`, etc.)

## Installation

### Option 1: Install from VSIX (Recommended)

1. Package the extension:
   ```bash
   npm install -g @vscode/vsce
   cd vscode-pace
   vsce package
   ```

2. Install in VS Code:
   - Press `Ctrl+Shift+P` (Windows/Linux) or `Cmd+Shift+P` (Mac)
   - Type "Extensions: Install from VSIX"
   - Select the generated `.vsix` file

### Option 2: Development Mode

1. Copy the `vscode-pace` folder to your VS Code extensions directory:
   - **Windows**: `%USERPROFILE%\.vscode\extensions\`
   - **macOS/Linux**: `~/.vscode/extensions/`

2. Restart VS Code

3. Open any `.pace` file to see syntax highlighting

## Usage

Simply open any file with a `.pace` extension, and syntax highlighting will be automatically applied.

## Example

```pace
# Define variables
set BUILD_OUTPUT "pace.exe"
set VERSION "1.0.0"

# Set default task
default "build"

# Create task alias
alias "b" "build"

# Define a task
task "build" {
    description "Build the Pace executable"
    command "go build -o ${BUILD_OUTPUT} cmd/pace/main.go"
    inputs ["cmd/**/*.go", "internal/**/*.go"]
    outputs ["${BUILD_OUTPUT}"]
    cache true
    timeout "5m"
}
```

## License

MIT
