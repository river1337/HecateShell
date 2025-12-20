#!/bin/bash
set -e

cd "$(dirname "$0")/hecate-shell-src"
echo "Building hecate..."
go build -o hecate .
echo "✓ Build complete: hecate-shell-src/hecate"
