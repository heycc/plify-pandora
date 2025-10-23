#!/bin/bash

# Build script for creating two separate WASM files using configurable architecture
# 1. official.wasm - Only official Go template functions
# 2. custom.wasm - Official functions + custom functions (getv, exists, get, jsonv)

echo "Building WASM files..."

# Build WASM with official Go template functions only
echo "Building official.wasm (official Go template functions only)..."
GOOS=js GOARCH=wasm go build -o official.wasm -tags "official" .

if [ $? -eq 0 ]; then
    echo "✓ official.wasm built successfully"
else
    echo "✗ Failed to build official.wasm"
    exit 1
fi

# Build WASM with custom functions
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
echo "Files created:"
echo "  - official.wasm (official Go template functions only)"
echo "  - custom.wasm (with custom functions: getv, exists, get, jsonv)"

# Show file sizes
echo ""
echo "File sizes:"
ls -lh *.wasm