#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VERSION="$(cat "$SCRIPT_DIR/../VERSION")"
OUT_NAME="tarmo-v$VERSION.exe"

"$SCRIPT_DIR/prod-build.sh"

WIN_USER="$(cmd.exe /c "echo %USERNAME%" 2>/dev/null | tr -d '\r')"
WIN_INSTALL_DIR="/mnt/c/Users/$WIN_USER/AppData/Local/Tarmo/$VERSION"
WIN_DATA_DIR="/mnt/c/ProgramData/Tarmo/data"

echo "Installing v$VERSION at $WIN_INSTALL_DIR"
mkdir -p "$WIN_INSTALL_DIR"
mkdir -p "$WIN_DATA_DIR"
cp -f "$SCRIPT_DIR/../build/$OUT_NAME" "$WIN_INSTALL_DIR/"

echo "Done:"
echo "  Executable -> $WIN_INSTALL_DIR/$OUT_NAME"
echo "  Data       -> $WIN_DATA_DIR"