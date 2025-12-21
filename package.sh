#!/bin/bash
set -e

# Package HecateShell for release
# Creates a minimal archive without dev files

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
VERSION=$(cat "$SCRIPT_DIR/version" 2>/dev/null || echo "dev")

# Output name must match what the installer expects to download
OUTPUT_NAME="HecateShell.tar.gz"

echo "Packaging HecateShell v${VERSION}..."

# Files/directories to EXCLUDE from the release
EXCLUDES=(
    "hecate-shell-src"
    ".git"
    ".gitignore"
    "*.md"
    "build.sh"
    "package.sh"
    "HecateShell.tar.gz"
)

# Build exclude args for tar
EXCLUDE_ARGS=""
for pattern in "${EXCLUDES[@]}"; do
    EXCLUDE_ARGS="$EXCLUDE_ARGS --exclude=$pattern"
done

# Create the archive
cd "$SCRIPT_DIR"
tar $EXCLUDE_ARGS -czvf "$OUTPUT_NAME" .

echo ""
echo "Created: $OUTPUT_NAME (v${VERSION})"
echo "Size: $(du -h "$OUTPUT_NAME" | cut -f1)"
echo ""
echo "Contents:"
tar -tzvf "$OUTPUT_NAME" | head -20
COUNT=$(tar -tzvf "$OUTPUT_NAME" | wc -l)
if [ "$COUNT" -gt 20 ]; then
    echo "... and $((COUNT - 20)) more files"
fi
echo ""
echo "Upload this file to GitHub releases as: HecateShell.tar.gz"
