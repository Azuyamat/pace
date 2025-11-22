#!/bin/sh
set -e

# --- Configuration ---
REPO="Azuyamat/pace"
BINARY="pace"
# ---------------------

echo "Initializing installer for $REPO..."

# 1. Detect OS & Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map hardware architecture to GoReleaser naming convention
case $ARCH in
    x86_64) ARCH="amd64";;
    aarch64) ARCH="arm64";;
    arm64) ARCH="arm64";; # macOS M1/M2
    i386) ARCH="386";;
    *) echo "Error: Unsupported architecture $ARCH"; exit 1;;
esac

# 2. Find Install Directory
# Try /usr/local/bin first (standard location). 
# If not writable and no sudo, fallback to ~/.local/bin
INSTALL_DIR="/usr/local/bin"
USE_SUDO=""

if [ ! -w "$INSTALL_DIR" ]; then
    if command -v sudo >/dev/null; then
        USE_SUDO="sudo"
    else
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
    fi
fi

# 3. Fetch Latest Release URL
API_URL="https://api.github.com/repos/$REPO/releases/latest"
echo "Fetching release info for $OS/$ARCH..."

# We use basic grep/cut to parse JSON to avoid depending on 'jq'
# Logic: Find 'browser_download_url', filter for OS+Arch+tar.gz
DOWNLOAD_URL=$(curl -sL "$API_URL" | grep "browser_download_url" | grep -i "$OS" | grep -i "$ARCH" | grep "tar.gz" | cut -d '"' -f 4 | head -n 1)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "Error: Could not find a release asset for $OS $ARCH."
    echo "Check if a release exists at https://github.com/$REPO/releases"
    exit 1
fi

echo "Downloading $(basename "$DOWNLOAD_URL")..."

# 4. Download & Extract
TMP_DIR=$(mktemp -d)
FILE_PATH="$TMP_DIR/download.tar.gz"

if curl -sL -o "$FILE_PATH" "$DOWNLOAD_URL"; then
    echo "Extracting..."
    tar -xzf "$FILE_PATH" -C "$TMP_DIR"
    
    # Find the binary in the extracted files (handles subfolders if they exist)
    BINARY_PATH=$(find "$TMP_DIR" -type f -name "$BINARY" | head -n 1)

    if [ -z "$BINARY_PATH" ]; then
        echo "Error: Could not find binary named '$BINARY' in archive."
        exit 1
    fi

    # 5. Install
    echo "Installing to $INSTALL_DIR..."
    $USE_SUDO mv "$BINARY_PATH" "$INSTALL_DIR/$BINARY"
    $USE_SUDO chmod +x "$INSTALL_DIR/$BINARY"
    
    # Cleanup
    rm -rf "$TMP_DIR"

    # 6. Path Check
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        echo ""
        echo "WARNING: $INSTALL_DIR is not in your PATH."
        echo "Add the following to your shell config (.bashrc / .zshrc):"
        echo "export PATH=\$PATH:$INSTALL_DIR"
    fi

    echo ""
    echo "Success! Run '$BINARY' to start."
else
    echo "Download failed."
    rm -rf "$TMP_DIR"
    exit 1
fi