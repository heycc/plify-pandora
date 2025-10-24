#!/bin/bash

# Build script for creating two separate WASM files using build tags
# 1. official.wasm - Only official Go template functions
# 2. custom.wasm - Official functions + custom functions (getv, exists, get, json)
#
# Architecture:
# - Core functionality is shared between both builds
# - Custom functions are conditionally compiled using build tags
# - functions_custom.go: included when NOT building with "official" tag
# - functions_official.go: included when building with "official" tag

set -e  # Exit on error

echo "Building WASM files with new architecture..."
echo ""

# Build WASM with official Go template functions only
echo "Building official.wasm (official Go template functions only)..."
GOOS=js GOARCH=wasm go build -tags official -o official.wasm .

if [ $? -eq 0 ]; then
    echo "✓ official.wasm built successfully"
else
    echo "✗ Failed to build official.wasm"
    exit 1
fi

echo ""

# Build WASM with custom functions (default build, no tags needed)
echo "Building custom.wasm (with custom functions)..."
GOOS=js GOARCH=wasm go build -o custom.wasm .

if [ $? -eq 0 ]; then
    echo "✓ custom.wasm built successfully"
else
    echo "✗ Failed to build custom.wasm"
    exit 1
fi

echo ""
echo "Build completed successfully!"
echo ""
echo "Files created:"
echo "  - official.wasm (official Go template functions only)"
echo "  - custom.wasm (with custom functions: getv, exists, get, json)"

# Show file sizes
echo ""
echo "File sizes:"
ls -lh official.wasm custom.wasm 2>/dev/null || ls -lh *.wasm

