#!/usr/bin/env bash
set -euo pipefail

workspace="${WORKSPACE_FOLDER:-$(pwd)}"
vsix_dir="$workspace/vscode-pace"

vsix="$(ls -t "$vsix_dir"/*.vsix 2>/dev/null | head -n 1 || true)"
if [[ -z "$vsix" ]]; then
  echo "No .vsix found in $vsix_dir" >&2
  exit 1
fi

exec_path="${VSCODE_EXEC_PATH:-}"

is_wsl=false
if grep -qi microsoft /proc/version 2>/dev/null; then
  is_wsl=true
fi

# In WSL/Remote, do NOT call the Windows Code.exe from ${execPath}; that can
# open a new local window. The `code` on PATH inside the VS Code integrated
# terminal is the Remote CLI and installs into the active remote/server side.
if command -v code >/dev/null 2>&1; then
  code_path="$(command -v code)"
  if [[ "$is_wsl" == true || "$code_path" == *vscode-server*remote-cli* || "$code_path" == *vscode-remote* ]]; then
    echo "Installing $vsix via VS Code Remote CLI: $code_path"
    code --install-extension "$vsix" --force
    exit $?
  fi
fi

# Native Linux/macOS VS Code: ${execPath} is directly executable.
if [[ -n "$exec_path" && -x "$exec_path" ]]; then
  echo "Installing $vsix via VS Code executable: $exec_path"
  "$exec_path" --install-extension "$vsix" --force
  exit $?
fi

# Non-remote fallback to the CLI on PATH.
if command -v code >/dev/null 2>&1; then
  echo "Installing $vsix via VS Code CLI: $(command -v code)"
  code --install-extension "$vsix" --force
  exit $?
fi

if command -v code-insiders >/dev/null 2>&1; then
  echo "Installing $vsix via VS Code Insiders CLI: $(command -v code-insiders)"
  code-insiders --install-extension "$vsix" --force
  exit $?
fi

echo "Could not find a usable VS Code CLI. VSCODE_EXEC_PATH='$exec_path'" >&2
exit 1
