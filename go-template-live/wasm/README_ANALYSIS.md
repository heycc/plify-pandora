# WASM Build Dependency Analysis Tools

This directory contains tools to analyze the dependencies for the three WASM build variants.

## Quick Start

```bash
# Visual overview (recommended for quick understanding)
./visualize_deps.sh

# Comprehensive Python analysis
python3 analyze_deps.py

# Detailed shell-based analysis
./analyze_dependencies.sh

# Read full documentation
cat DEPENDENCIES.md
```

## Build Variants

The `build.sh` script creates three WASM files using Go build tags:

### 1. official.wasm
- **Tag:** `official`
- **Size:** ~4.70 MB
- **Functions:** Only standard Go template functions
- **Files:** 7 (5 core + 2 specific)
- **Use case:** When you only need standard Go templates

### 2. custom.wasm
- **Tag:** `custom`
- **Size:** ~4.72 MB
- **Functions:** Standard + 5 custom functions (getv, exists, get, json, jsonArray)
- **Files:** 7 (5 core + 2 specific)
- **Use case:** When you need simple custom helper functions

### 3. confd.wasm (Default)
- **Tag:** `confd`
- **Size:** ~4.78 MB
- **Functions:** Standard + 25+ Confd-style functions
- **Files:** 7 (5 core + 2 specific)
- **Use case:** Full-featured template processing (default)

## Analysis Tools

### 1. visualize_deps.sh
**Best for:** Quick visual overview

```bash
./visualize_deps.sh
```

Shows:
- Core files included in all builds
- Build-specific files for each variant
- ASCII art dependency flow diagram
- File sizes and function counts
- Quick reference commands

### 2. analyze_deps.py
**Best for:** Comprehensive analysis

```bash
python3 analyze_deps.py
```

Shows:
- Complete file breakdown by build
- Import analysis per build
- Build composition details
- Test file identification
- File sizes and statistics

### 3. analyze_dependencies.sh
**Best for:** Shell-based detailed analysis

```bash
./analyze_dependencies.sh
```

Shows:
- Build constraints for each file
- Import analysis by file
- Dependency graph
- Function availability by build
- Summary and tips

### 4. DEPENDENCIES.md
**Best for:** Complete documentation

```bash
cat DEPENDENCIES.md
# or
open DEPENDENCIES.md
```

Contains:
- Full architecture explanation
- Detailed file breakdown
- Import analysis
- Build size comparison
- How-to guides
- Key insights

## Core Architecture

### Core Files (Always Included)
These 5 files are included in **all** WASM builds:

```
main.go              - Entry point, WASM initialization
parser.go            - Template parsing & variable extraction
wasm_handlers.go     - JavaScript/WASM interface
types.go             - Type definitions
function_base.go     - Helper functions
```

### Build-Specific Files (Conditionally Included)

Each build includes **one** of these function files based on the build tag:

```
functions_official.go    - Tag: js && official
functions_custom.go      - Tag: js && custom
functions_confd.go       - Tag: js && confd
functions_default.go     - Tag: js && !custom && !confd && !official
```

## Dependency Flow

```
Core Files (5)
    │
    ├─> official tag ─> functions_official.go ─> official.wasm
    ├─> custom tag   ─> functions_custom.go   ─> custom.wasm
    └─> confd tag    ─> functions_confd.go    ─> confd.wasm (→ main.wasm)
```

## Go Commands for Analysis

### List files in a build
```bash
GOOS=js GOARCH=wasm go list -tags official -f '{{range .GoFiles}}{{println .}}{{end}}' .
GOOS=js GOARCH=wasm go list -tags custom -f '{{range .GoFiles}}{{println .}}{{end}}' .
GOOS=js GOARCH=wasm go list -tags confd -f '{{range .GoFiles}}{{println .}}{{end}}' .
```

### Show all dependencies
```bash
GOOS=js GOARCH=wasm go list -tags official -deps .
GOOS=js GOARCH=wasm go list -tags custom -deps .
GOOS=js GOARCH=wasm go list -tags confd -deps .
```

### Why is a package included?
```bash
GOOS=js GOARCH=wasm go mod why -tags confd encoding/base64
```

### Show import graph
```bash
GOOS=js GOARCH=wasm go list -tags confd -f '{{.ImportPath}} -> {{join .Imports ", "}}' .
```

## Key Insights

1. **Shared Core**: All builds share 5 core files (~90% of code)
2. **Minimal Size Impact**: Custom functions add only ~2% to WASM size
3. **Clean Separation**: Build tags ensure no function conflicts
4. **No External Deps**: Only Go standard library is used
5. **Test Isolation**: Test files properly excluded with `!js` tag

## File Statistics

| Build | Total Files | Core Files | Specific Files | Size | Functions |
|-------|-------------|------------|----------------|------|-----------|
| official.wasm | 7 | 5 | 2 | 4.70 MB | 0 custom |
| custom.wasm | 7 | 5 | 2 | 4.72 MB | 5 custom |
| confd.wasm | 7 | 5 | 2 | 4.78 MB | 25+ custom |

## Import Comparison

### All Builds
```
encoding/json
errors
fmt
strings
syscall/js
text/template
text/template/parse
```

### confd.wasm Additional Imports
```
encoding/base64
path
strconv
time
```

## Build Process

The `build.sh` script:
1. Builds `official.wasm` with `-tags official`
2. Builds `custom.wasm` with `-tags custom`
3. Builds `confd.wasm` with `-tags confd`
4. Copies `confd.wasm` to `main.wasm` (default)
5. Copies `main.wasm` to frontend public directory

## Troubleshooting

### Build fails
```bash
# Check which files would be included
GOOS=js GOARCH=wasm go list -tags confd -f '{{.GoFiles}}'

# Check for build constraint errors
go vet -tags confd .
```

### Wrong functions available
```bash
# Verify which WASM is loaded
ls -lh *.wasm

# Check which functions are registered
grep -A 20 "CreateRenderFuncMap" functions_*.go
```

### Size concerns
```bash
# Compare sizes
ls -lh official.wasm custom.wasm confd.wasm

# Analyze what's taking space
go tool nm confd.wasm | sort -k2 -r | head -20
```

## Contributing

When adding new functions:

1. Choose the appropriate file:
   - `functions_custom.go` for simple helpers
   - `functions_confd.go` for Confd-style functions
   - Create new file with appropriate build tags for new categories

2. Add build tags at the top:
   ```go
   //go:build js && yourtag
   // +build js,yourtag
   ```

3. Update `build.sh` to build your variant

4. Run analysis tools to verify:
   ```bash
   ./visualize_deps.sh
   python3 analyze_deps.py
   ```

## References

- [Go Build Constraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints)
- [WebAssembly in Go](https://github.com/golang/go/wiki/WebAssembly)
- [Confd Template Functions](https://github.com/kelseyhightower/confd/blob/master/docs/templates.md)

