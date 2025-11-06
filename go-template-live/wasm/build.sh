#!/bin/bash

# Build script for creating three separate WASM files using build tags
# 1. official.wasm - Only official Go template functions
# 2. custom.wasm - Official functions + custom functions (getv, exists, get, json)
# 3. confd.wasm - Official functions + Confd-style functions
#
# Architecture:
# - Core functionality is shared between all builds
# - Custom functions are conditionally compiled using build tags
# - functions_custom.go: included when building with "custom" tag
# - functions_confd.go: included when building with "confd" tag
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

# Build WASM with custom functions
echo "Building custom.wasm (with custom functions)..."
GOOS=js GOARCH=wasm go build -tags custom -o custom.wasm .

if [ $? -eq 0 ]; then
    echo "✓ custom.wasm built successfully"
else
    echo "✗ Failed to build custom.wasm"
    exit 1
fi

echo ""

# Build WASM with Confd-style functions
echo "Building confd.wasm (with Confd-style functions)..."
GOOS=js GOARCH=wasm go build -tags confd -o confd.wasm .

if [ $? -eq 0 ]; then
    echo "✓ confd.wasm built successfully"
else
    echo "✗ Failed to build confd.wasm"
    exit 1
fi

echo ""

# Copy confd.wasm to main.wasm as the default WASM for frontend
echo "Copying confd.wasm to main.wasm (default WASM for frontend)..."
cp confd.wasm main.wasm
echo "✓ main.wasm created"

# Copy to frontend public directory if it exists
FRONTEND_PUBLIC="../frontend/public"
if [ -d "$FRONTEND_PUBLIC" ]; then
    echo "Copying main.wasm to $FRONTEND_PUBLIC..."
    cp main.wasm "$FRONTEND_PUBLIC/"
    echo "✓ main.wasm copied to frontend"
else
    echo "⚠ Frontend public directory not found at $FRONTEND_PUBLIC"
fi

echo ""
echo "Build completed successfully!"
echo ""
echo "Files created:"
echo "  - official.wasm (official Go template functions only)"
echo "  - custom.wasm (with custom functions: getv, exists, get, json, jsonArray)"
echo "  - confd.wasm (with Confd-style functions: base, split, json, jsonArray, dir, map, join, datetime, toUpper, toLower, replace, contains, base64Encode, base64Decode, trimSuffix, parseBool, reverse, add, sub, div, mod, mul, seq, atoi)"
echo "  - main.wasm (copy of confd.wasm for frontend)"

# Show file sizes
echo ""
echo "File sizes:"
ls -lh official.wasm custom.wasm confd.wasm main.wasm 2>/dev/null || ls -lh *.wasm

