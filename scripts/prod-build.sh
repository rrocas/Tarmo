#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VERSION="$(cat "$SCRIPT_DIR/../VERSION")"
OUT_NAME="tarmo-v$VERSION.exe"
MAIN_DIR="$SCRIPT_DIR/../cmd/rest_api"

echo "Generating Windows icon resource..."
x86_64-w64-mingw32-windres "$MAIN_DIR/app.rc" -O coff -o "$MAIN_DIR/rsrc_windows_amd64.syso"

echo "Building Tarmo v$VERSION for Windows (cgo cross-compile, static link)..."

CGO_ENABLED=1 \
GOOS=windows \
GOARCH=amd64 \
CC=x86_64-w64-mingw32-gcc \
CGO_LDFLAGS="-static" \
go build -ldflags "-X tarmo/internal/config.Env=production -X tarmo/internal/config.Version=$VERSION -linkmode external -extldflags -static" \
  -o "build/$OUT_NAME" "$MAIN_DIR"

echo "Build ready: $(pwd)/build/$OUT_NAME"