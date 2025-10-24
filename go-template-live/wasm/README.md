# Go Template Live - WASM Module

This project provides a WebAssembly (WASM) module for parsing Go templates and extracting variables. It supports two build configurations:

- **Official Functions Only**: Standard Go template functions
- **Custom Functions**: Official functions plus Confd-like custom functions (getv, exists, get, json)

## 🚀 Quick Start

### Build Both WASM Files

```bash
# Use the build script to create both versions
./build.sh

# This creates:
# - official.wasm (official Go template functions only)
# - custom.wasm (official functions + custom functions)
```

## 🏗️ Architecture Overview

The project uses a **plugin-style architecture** with **build tags** for clean separation between official and custom builds.

### Key Design Principles

1. **Build Tags for Compile-Time Selection**: Uses Go build tags to include/exclude custom functions at compile time
2. **Single Source of Truth**: `FunctionRegistry` is the only place that tracks available functions
3. **Easy Extensibility**: Add new functions by creating new files with appropriate build tags
4. **Code Reuse**: Core parsing logic is shared; only function implementations differ

### Core Components

#### Type Definitions (`types.go`)
- `VariableInfo`: Stores variable name and optional default value
- `FunctionDefinition`: Defines function name, handler, and variable extractors
- `FunctionRegistry`: Central registry for all template functions

#### Function Base (`function_base.go`)
- Helper functions for extracting variables from function arguments
- Global registry instance (`globalRegistry`)
- No build tags - included in both builds

#### Custom Functions (`functions_custom.go`)
- **Build Tag**: `!official` (included when NOT building with "official" tag)
- Registers custom functions via `init()` function
- Implements: `getv`, `exists`, `get`, `json`
- Inspired by [Confd](https://github.com/kelseyhightower/confd)

#### Official Build (`functions_official.go`)
- **Build Tag**: `official` (included only when building with "official" tag)
- Provides empty `CreateRenderFuncMap()` for official builds
- No custom functions registered

#### Parser (`parser.go`)
- Extracts variables from template AST
- Uses `FunctionRegistry` to determine which functions to recognize
- Supports both simple extraction and extraction with default values

#### WASM Handler (`wasm_handler.go`)
- JavaScript interface for WASM module
- Functions exposed: `extractTemplateVariables`, `extractTemplateVariablesSimple`, `renderTemplateWithValues`

### Custom Functions

| Function | Description | Example |
|----------|-------------|---------|
| `getv` | Get variable with optional default | `{{getv "username" "guest"}}` |
| `exists` | Check if variable exists | `{{exists "feature_flag"}}` |
| `get` | Get variable (errors if not found) | `{{get "required_field"}}` |
| `json` | Parse JSON variable | `{{json "config_json"}}` |

## 📦 Build Process

### Prerequisites
- Go 1.16+
- WASM target support

### Building WASM Files

Use the provided build script:

```bash
./build.sh
```

This creates two WASM files:
- `official.wasm` - Official Go template functions only
- `custom.wasm` - Official functions + custom functions

### Manual Build Commands

```bash
# Build official version (no custom functions)
GOOS=js GOARCH=wasm go build -tags official -o official.wasm .

# Build custom version (with custom functions)
GOOS=js GOARCH=wasm go build -o custom.wasm .
```

### Build Tags Explained

- **`official`**: When present, excludes custom functions
  - Includes: `functions_official.go`
  - Excludes: `functions_custom.go`
  
- **No tags (default)**: Includes custom functions
  - Includes: `functions_custom.go`
  - Excludes: `functions_official.go`

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test -v

# Run tests with coverage
go test -cover

# Run specific test categories
go test -v -run "_Pure"        # Pure tests only
go test -v -run "TestParser"   # Parser tests only
```

### Test Structure

#### Test Files
- `template_parser_pure_test.go` - Unit tests for individual components
- `template_parser_test.go` - Integration tests for both configurations
- `test_helpers.go` - Test utilities for dual-configuration testing

#### Test Helper (`test_helpers.go`)
Provides utilities for testing both configurations:

```go
helper := NewTestHelper()

// Test custom functions configuration
parserCustom := helper.NewParserWithCustomFunctions()

// Test official functions configuration
parserOfficial := helper.NewParserWithOfficialFunctions()
```

### Test Coverage

Tests verify:
- Variable extraction from standard template fields
- Variable extraction from custom function parameters
- Both build configurations work correctly
- Function registry operations
- Template rendering with variables

## 📖 Usage

### JavaScript Interface

After loading the WASM module, these functions are available:

```javascript
// Extract variables with default values
const variables = extractTemplateVariables(templateContent, fileName);

// Extract only variable names (no defaults)
const variableNames = extractTemplateVariablesSimple(templateContent, fileName);

// Render template with variable values
const result = renderTemplateWithValues(templateContent, variablesJSON);
```

### Example Usage

```javascript
// Extract variables from template
const template = `Hello {{.Name}}, your username is {{getv "username" "guest"}}`;
const variables = JSON.parse(extractTemplateVariables(template));
// variables = [{name: "Name"}, {name: "username", defaultValue: "guest"}]

// Render template
const values = { Name: "John", username: "john_doe" };
const result = renderTemplateWithValues(template, JSON.stringify(values));
// result = "Hello John, your username is john_doe"
```

## 📁 File Structure

```
wasm/
├── types.go                          # Core types and registry
├── function_base.go                  # Base function utilities
├── functions_custom.go               # Custom functions (build tag: !official)
├── functions_official.go             # Official build stub (build tag: official)
├── parser.go                     # Template parser
├── wasm_handler.go               # WASM/JavaScript interface
├── main.go                           # WASM entry point
├── build.sh                      # Build script
├── test_helpers.go               # Test utilities
├── template_parser_pure_test.go  # Unit tests
├── template_parser_test.go           # Integration tests
└── README.md                     # This file
```

## 🔧 Extending with New Functions

The new architecture makes it easy to add custom functions:

### 1. Add Function Implementation

Create a new file with the `!official` build tag (or add to `functions_custom.go`):

```go
//go:build !official
// +build !official

package main

// Add to init() function in functions_custom.go
func init() {
    registry := GetGlobalRegistry()
    
    registry.RegisterFunction(&FunctionDefinition{
        Name:                  "newfunc",
        Description:           "New function description",
        Handler:               newfuncMinimalHandler,
        Extractor:             extractNewfuncVariables,
        ExtractorWithDefaults: extractNewfuncVariablesWithDefaults,
    })
}

// Minimal handler for parsing
func newfuncMinimalHandler(args ...string) string {
    return ""
}

// Render handler with actual logic
func newfuncRenderHandler(variables map[string]interface{}) func(args ...string) string {
    return func(args ...string) string {
        // Implementation
        return ""
    }
}

// Variable extractors
func extractNewfuncVariables(args []parse.Node, cycle int) ([]string, error) {
    // Extract variable names from arguments
    return extractStringArgVariable(args, cycle, 1)
}

func extractNewfuncVariablesWithDefaults(args []parse.Node, cycle int) ([]VariableInfo, error) {
    // Extract variables with defaults
    return extractStringArgVariableWithDefaults(args, cycle, 1, 2)
}
```

### 2. Add to Render Function Map

Update `CreateRenderFuncMap()` in `functions_custom.go`:

```go
func CreateRenderFuncMap(variables map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "getv":    getvRenderHandler(variables),
        "exists":  existsRenderHandler(variables),
        "get":     getRenderHandler(variables),
        "json":   jsonRenderHandler(variables),
        "newfunc": newfuncRenderHandler(variables),  // Add your function
    }
}
```

### 3. Rebuild

```bash
./build.sh
```

That's it! Your new function is now available in the custom build.

## 🎯 Benefits of New Architecture

✅ **Cleaner Separation**: Build tags handle official vs custom at compile time  
✅ **Easy to Extend**: Add new functions by creating new files with build tags  
✅ **Less Code Duplication**: Shared core logic, only functions differ  
✅ **Better Maintainability**: Clear structure, single responsibility per file  
✅ **Confd-style Extensibility**: Easy to add more custom functions like Confd  
✅ **No Runtime Configuration**: All decisions made at compile time for better performance  
✅ **Single Source of Truth**: FunctionRegistry is the only place tracking functions  

## 🤝 Contributing

Feel free to:
- Add more template examples
- Improve the architecture
- Add new custom functions
- Report bugs
- Suggest improvements
