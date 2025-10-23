# Go Template Live - WASM Module

This project provides a WebAssembly (WASM) module for parsing Go templates and extracting variables. It supports two build configurations:

- **Official Functions Only**: Standard Go template functions
- **Custom Functions**: Official functions plus Confd-like custom functions (getv, exists, get, jsonv)

## 🚀 Quick Start

### Build Both WASM Files

```bash
# Use the build script to create both versions
./build.sh

# This creates:
# - official.wasm (official Go template functions only)
# - custom.wasm (official functions + custom functions)
```

## Architecture Overview

### Core Components

#### Build Configuration (`build_config.go`)
- `BuildConfig` struct controls which functions are included
- `DefaultBuildConfig()` - Includes custom functions
- `OfficialBuildConfig()` - Official functions only

#### Function Registry (`functions.go`)
- `FunctionRegistry` manages custom template functions
- `FunctionDefinition` defines function name, handler, and variable extractor
- Always registers minimal implementations for parsing, execution behavior differs by config

#### Function Matcher (`matcher.go`)
- `FunctionMatcher` interface for matching custom functions
- `DefaultFunctionMatcher` - Recognizes custom functions
- `OfficialFunctionMatcher` - Only recognizes official functions

#### Template Parser (`template_parser.go`, `parser.go`)
- `Parser` extracts variables from templates
- Supports both standard template fields and custom function parameters
- Configurable via `BuildConfig` and `FunctionMatcher`

#### WASM Interface (`wasm_handlers.go`)
- `WASMHandler` provides JavaScript interface
- Functions exposed: `extractTemplateVariables`, `extractTemplateVariablesSimple`, `renderTemplateWithValues`

### Custom Functions

| Function | Description | Example |
|----------|-------------|---------|
| `getv` | Get variable with optional default | `{{getv "username" "guest"}}` |
| `exists` | Check if variable exists | `{{exists "feature_flag"}}` |
| `get` | Get variable (errors if not found) | `{{get "required_field"}}` |
| `jsonv` | Parse JSON variable | `{{jsonv "config_json"}}` |

## Build Process

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
GOOS=js GOARCH=wasm go build -o official.wasm .

# Build custom version (with custom functions)
GOOS=js GOARCH=wasm go build -tags custom -o custom.wasm .
```

### Build Tags

- `js` - Required for WASM builds
- `custom` - Includes custom functions (used in custom.wasm)

## Testing

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
- Function matching logic
- Template rendering with variables

## Usage

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
├── build_config.go      # Build configuration
├── functions.go         # Function registry and definitions
├── matcher.go           # Function matching logic
├── parser.go            # Template parser
├── template_parser.go   # Template parsing utilities
├── wasm_handlers.go     # WASM/JavaScript interface
├── main.go              # WASM entry point
├── build.sh             # Build script
├── README.md            # This file
└── *_test.go           # Test files
```

## Extending with New Functions

To add new custom functions:

1. **Add to Function Registry** (`functions.go`):
   ```go
   registry.RegisterFunction(&FunctionDefinition{
       Name:        "newfunc",
       Description: "New function description",
       Handler:     newFuncHandler,
       Extractor:   extractNewFuncVariables,
   })
   ```

2. **Add to Function Matcher** (`matcher.go`):
   ```go
   supportedFunctions: map[string]bool{
       // ... existing functions
       "newfunc": true,
   }
   ```

3. **Implement Handler and Extractor**:
   ```go
   func newFuncHandler(variables map[string]interface{}) func(args ...string) string {
       return func(args ...string) string {
           // Implementation
       }
   }

   func extractNewFuncVariables(args []parse.Node, cycle int) ([]string, error) {
       // Extract variable names from arguments
   }
   ```

## 🤝 Contributing

Feel free to:
- Add more template examples
- Improve the architecture
- Add new custom functions
- Report bugs
- Suggest improvements