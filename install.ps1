<#
.SYNOPSIS
    Installs a binary from GitHub Releases (GoReleaser compatible).
    
.DESCRIPTION
    Fetches the latest release from GitHub, matches the asset based on 
    OS (Windows) and Architecture, extracts the binary, and adds it to PATH.
#>

param(
    [string]$Repo = "github.com/Azuyamat/pace",
    [string]$BinaryName = "pace"
)

$ErrorActionPreference = "Stop"

# --- Helper: Get Latest Release & Asset URL ---
function Get-DownloadUrl {
    param([string]$Repository, [string]$Arch)
    
    # Clean repo string
    $ghRepo = $Repository -replace "github.com/", ""
    $apiUrl = "https://api.github.com/repos/$ghRepo/releases/latest"
    
    try {
        Write-Host "Fetching release info..." -ForegroundColor Gray
        $release = Invoke-RestMethod -Uri $apiUrl
    } catch {
        throw "Could not fetch releases for $ghRepo. Is it public?"
    }

    # GoReleaser Arch Mapping (Adjust these patterns if your file naming differs)
    # We look for 'windows' AND the architecture in the filename
    $matchString = ""
    if ($Arch -eq "amd64") { $matchString = "windows.*(amd64|x86_64|64bit)" }
    elseif ($Arch -eq "arm64") { $matchString = "windows.*(arm64)" }
    elseif ($Arch -eq "386")   { $matchString = "windows.*(i386|386|32bit)" }

    # Filter assets
    $asset = $release.assets | Where-Object { $_.name -match "(?i)$matchString" -and $_.name -match "(?i)\.(zip|tar\.gz)$" } | Select-Object -First 1
    
    if (-not $asset) {
        throw "No compatible Windows binary found in release $($release.tag_name) for architecture $Arch."
    }

    Write-Host "Found version: $($release.tag_name)" -ForegroundColor Cyan
    return @{ Url = $asset.browser_download_url; Name = $asset.name; Version = $release.tag_name }
}

# --- Step 1: Detect System ---
Write-Host "Initializing..." -ForegroundColor Cyan

# Architecture Detection
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { $sysArch = "arm64" }
elseif ($env:PROCESSOR_ARCHITECTURE -eq "x86") { $sysArch = "386" }
else { $sysArch = "amd64" }

# Admin Detection
$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if ($isAdmin) {
    $installDir = "$env:ProgramFiles\$BinaryName"
    $pathScope = "Machine"
} else {
    $installDir = "$env:LOCALAPPDATA\Programs\$BinaryName"
    $pathScope = "User"
}

# --- Step 2: Find Asset ---
try {
    $downloadInfo = Get-DownloadUrl -Repository $Repo -Arch $sysArch
} catch {
    Write-Error $_
    exit 1
}

# --- Step 3: Download & Extract ---
$tempDir = Join-Path $env:TEMP "pace_install"
if (Test-Path $tempDir) { Remove-Item -Recurse -Force $tempDir }
New-Item -ItemType Directory -Force -Path $tempDir | Out-Null

$zipPath = Join-Path $tempDir $downloadInfo.Name
$exeDest = Join-Path $installDir "$BinaryName.exe"

Write-Host "Downloading $($downloadInfo.Name)..." -ForegroundColor Cyan
Invoke-WebRequest -Uri $downloadInfo.Url -OutFile $zipPath

Write-Host "Extracting..." -ForegroundColor Cyan
# Handle ZIP or TAR.GZ
if ($zipPath -match "\.zip$") {
    Expand-Archive -Path $zipPath -DestinationPath $tempDir -Force
} elseif ($zipPath -match "\.tar\.gz$") {
    # Windows 10/11 includes tar built-in
    tar -xf $zipPath -C $tempDir
}

# Find the .exe inside the extracted folder (it might be in a subfolder)
$foundBinary = Get-ChildItem -Path $tempDir -Recurse -Filter "$BinaryName.exe" | Select-Object -First 1

if (-not $foundBinary) {
    Write-Error "Could not find '$BinaryName.exe' inside the downloaded archive."
    exit 1
}

# --- Step 4: Install ---
if (-not (Test-Path $installDir)) { New-Item -ItemType Directory -Force -Path $installDir | Out-Null }

Write-Host "Installing to $installDir..." -ForegroundColor Cyan
if (Test-Path $exeDest) {
    $oldExe = Join-Path $installDir "$BinaryName.exe.old"
    if (Test-Path $oldExe) {
        Remove-Item -Path $oldExe -Force -ErrorAction SilentlyContinue
    }
    Move-Item -Path $exeDest -Destination $oldExe -Force
}
Move-Item -Path $foundBinary.FullName -Destination $exeDest -Force

# --- Step 5: Update PATH ---
$currentPath = [Environment]::GetEnvironmentVariable("Path", $pathScope)
if ($currentPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$currentPath;$installDir", $pathScope)
    $env:PATH += ";$installDir"
    Write-Host "Added to PATH." -ForegroundColor Green
}

# Cleanup
Remove-Item -Recurse -Force $tempDir

Write-Host "`nSuccess! Run '$BinaryName' to start." -ForegroundColor Green