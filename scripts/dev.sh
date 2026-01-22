#!/bin/bash
# Development startup script for AutoBMAD
# Starts both the Golang backend and Electron frontend concurrently

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "Starting AutoBMAD development environment..."

# Build the core binary first
echo "Building Go backend..."
cd "$PROJECT_ROOT/apps/core"
go build -o autobmad ./cmd/autobmad

echo "Go backend built successfully."
echo ""
echo "Starting Electron app with hot reload..."
echo "The Go backend will be spawned by Electron main process in production."
echo ""

# Start the Electron development server
cd "$PROJECT_ROOT"
pnpm run dev
