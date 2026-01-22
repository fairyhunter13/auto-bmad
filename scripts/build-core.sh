#!/bin/bash
# Build script for AutoBMAD Core (Golang backend)
# Supports cross-compilation for Linux and macOS

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
CORE_DIR="$PROJECT_ROOT/apps/core"
OUTPUT_DIR="$PROJECT_ROOT/apps/desktop/resources/bin"

# Version information
VERSION="${VERSION:-dev}"
COMMIT="${COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo 'none')}"
DATE="${DATE:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}"

LDFLAGS="-s -w -X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE"

echo "Building AutoBMAD Core..."
echo "  Version: $VERSION"
echo "  Commit:  $COMMIT"
echo "  Date:    $DATE"

mkdir -p "$OUTPUT_DIR"

cd "$CORE_DIR"

# Determine target platform(s)
if [ "$1" == "all" ]; then
    TARGETS=("linux-amd64" "linux-arm64" "darwin-amd64" "darwin-arm64")
elif [ -n "$1" ]; then
    TARGETS=("$1")
else
    # Default: build for current platform
    GOOS=$(go env GOOS)
    GOARCH=$(go env GOARCH)
    TARGETS=("$GOOS-$GOARCH")
fi

for target in "${TARGETS[@]}"; do
    GOOS="${target%-*}"
    GOARCH="${target#*-}"
    
    OUTPUT_NAME="autobmad-$GOOS-$GOARCH"
    if [ "$GOOS" == "windows" ]; then
        OUTPUT_NAME="$OUTPUT_NAME.exe"
    fi
    
    echo "Building for $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "$LDFLAGS" -o "$OUTPUT_DIR/$OUTPUT_NAME" ./cmd/autobmad
    echo "  Created: $OUTPUT_DIR/$OUTPUT_NAME"
done

echo "Build complete!"
