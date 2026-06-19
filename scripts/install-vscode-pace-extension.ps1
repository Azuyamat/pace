$ErrorActionPreference = 'Stop'

$workspace = if ($env:WORKSPACE_FOLDER) { $env:WORKSPACE_FOLDER } else { (Get-Location).Path }
$vsixDir = Join-Path $workspace 'vscode-pace'
$vsix = Get-ChildItem -Path $vsixDir -Filter '*.vsix' -File | Sort-Object LastWriteTime -Descending | Select-Object -First 1

if (-not $vsix) {
    throw "No .vsix found in $vsixDir"
}

$execPath = $env:VSCODE_EXEC_PATH

if ($execPath -and (Test-Path $execPath)) {
    & $execPath --install-extension $vsix.FullName --force
    exit $LASTEXITCODE
}

$code = Get-Command code -ErrorAction SilentlyContinue
if ($code) {
    & $code.Source --install-extension $vsix.FullName --force
    exit $LASTEXITCODE
}

$codeInsiders = Get-Command code-insiders -ErrorAction SilentlyContinue
if ($codeInsiders) {
    & $codeInsiders.Source --install-extension $vsix.FullName --force
    exit $LASTEXITCODE
}

throw "Could not find a usable VS Code executable. VSCODE_EXEC_PATH='$execPath'"
