#!/bin/bash

# Script to analyze dependencies for each WASM build tag
# This helps understand which files are included in each build variant

set -e

echo "=========================================="
echo "WASM Build Dependency Analysis"
echo "=========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to analyze dependencies for a specific build tag
analyze_build() {
    local tag=$1
    local output_file=$2
    
    echo -e "${BLUE}Analyzing dependencies for build tag: ${GREEN}${tag}${NC}"
    echo ""
    
    # Use go list to get dependencies
    echo "Files included in build:"
    GOOS=js GOARCH=wasm go list -tags "$tag" -f '{{range .GoFiles}}{{println .}}{{end}}' . | sort
    
    echo ""
    echo "Dependencies (packages):"
    GOOS=js GOARCH=wasm go list -tags "$tag" -f '{{range .Deps}}{{println .}}{{end}}' . | grep -v "internal" | sort | uniq
    
    echo ""
    echo "Build constraints for each .go file:"
    for file in *.go; do
        if [ -f "$file" ]; then
            # Extract build tags from file
            tags=$(head -n 5 "$file" | grep -E "^//go:build|^\+build" || true)
            if [ -n "$tags" ]; then
                echo -e "  ${YELLOW}$file${NC}:"
                echo "$tags" | sed 's/^/    /'
            else
                echo -e "  ${YELLOW}$file${NC}: (no build constraints - included in all builds)"
            fi
        fi
    done
    
    echo ""
    echo "----------------------------------------"
    echo ""
}

# Function to show file sizes and compare
show_sizes() {
    echo -e "${BLUE}Build Output Sizes:${NC}"
    echo ""
    
    if [ -f "official.wasm" ]; then
        size=$(ls -lh official.wasm | awk '{print $5}')
        echo -e "  official.wasm: ${GREEN}$size${NC}"
    fi
    
    if [ -f "custom.wasm" ]; then
        size=$(ls -lh custom.wasm | awk '{print $5}')
        echo -e "  custom.wasm:   ${GREEN}$size${NC}"
    fi
    
    if [ -f "confd.wasm" ]; then
        size=$(ls -lh confd.wasm | awk '{print $5}')
        echo -e "  confd.wasm:    ${GREEN}$size${NC}"
    fi
    
    echo ""
}

# Function to analyze imports in each file
analyze_imports() {
    echo -e "${BLUE}Import Analysis by File:${NC}"
    echo ""
    
    for file in *.go; do
        if [ -f "$file" ]; then
            # Skip test files
            if [[ "$file" == *_test.go ]]; then
                continue
            fi
            
            # Extract build tags
            tags=$(head -n 5 "$file" | grep -E "^//go:build" | sed 's|//go:build ||' || echo "all builds")
            
            # Extract imports
            imports=$(awk '/^import \(/,/^\)/' "$file" | grep -v "^import" | grep -v "^)" | sed 's/^[[:space:]]*//' | sed 's/"//g' || true)
            
            if [ -n "$imports" ]; then
                echo -e "${YELLOW}$file${NC} (${tags}):"
                echo "$imports" | sed 's/^/  /'
                echo ""
            fi
        fi
    done
}

# Function to create a dependency graph in text format
create_dependency_graph() {
    echo -e "${BLUE}Dependency Graph:${NC}"
    echo ""
    echo "Build Tag -> Files -> Core Dependencies"
    echo ""
    
    for tag in "official" "custom" "confd"; do
        echo -e "${GREEN}${tag}${NC}:"
        
        # Get Go files for this build
        files=$(GOOS=js GOARCH=wasm go list -tags "$tag" -f '{{range .GoFiles}}{{println .}}{{end}}' . 2>/dev/null | sort)
        
        echo "  Go Files:"
        echo "$files" | sed 's/^/    - /'
        
        echo ""
    done
}

# Function to show which functions are available in each build
analyze_functions() {
    echo -e "${BLUE}Function Availability by Build:${NC}"
    echo ""
    
    echo -e "${GREEN}official.wasm${NC}:"
    echo "  - Only standard Go template functions (no custom functions)"
    echo ""
    
    echo -e "${GREEN}custom.wasm${NC}:"
    echo "  - Standard Go template functions"
    echo "  - Custom functions from functions_custom.go:"
    if [ -f "functions_custom.go" ]; then
        grep -E "^\s*\"[a-zA-Z]+\":" functions_custom.go | sed 's/://g' | sed 's/"//g' | sed 's/^/    - /'
    fi
    echo ""
    
    echo -e "${GREEN}confd.wasm${NC}:"
    echo "  - Standard Go template functions"
    echo "  - Confd-style functions from functions_confd.go:"
    if [ -f "functions_confd.go" ]; then
        grep -E "^\s*\"[a-zA-Z0-9]+\":" functions_confd.go | sed 's/://g' | sed 's/"//g' | sed 's/^/    - /'
    fi
    echo ""
}

# Main analysis
echo ""

# Analyze each build tag
analyze_build "official" "official_deps.txt"
analyze_build "custom" "custom_deps.txt"
analyze_build "confd" "confd_deps.txt"

# Show file sizes if WASM files exist
if [ -f "official.wasm" ] || [ -f "custom.wasm" ] || [ -f "confd.wasm" ]; then
    show_sizes
fi

# Analyze imports
analyze_imports

# Create dependency graph
create_dependency_graph

# Analyze available functions
analyze_functions

# Summary
echo "=========================================="
echo "Summary"
echo "=========================================="
echo ""
echo "Core files (included in ALL builds):"
echo "  - main.go (entry point)"
echo "  - parser.go (template parsing logic)"
echo "  - wasm_handlers.go (JavaScript interface)"
echo "  - types.go (type definitions)"
echo "  - function_base.go (helper functions)"
echo ""
echo "Build-specific files (conditionally compiled):"
echo "  - functions_official.go (tag: official)"
echo "  - functions_custom.go (tag: custom)"
echo "  - functions_confd.go (tag: confd)"
echo ""
echo "Test files (not included in WASM builds):"
echo "  - *_test.go files"
echo ""
echo "To see detailed dependency tree for a specific build:"
echo "  GOOS=js GOARCH=wasm go list -tags official -deps ."
echo "  GOOS=js GOARCH=wasm go list -tags custom -deps ."
echo "  GOOS=js GOARCH=wasm go list -tags confd -deps ."
echo ""
echo "To see why a package is included:"
echo "  GOOS=js GOARCH=wasm go mod why -tags official <package>"
echo ""

