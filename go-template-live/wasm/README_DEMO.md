# Go Template Parser WASM Demo

A comprehensive interactive demo showcasing the Go template parser's WebAssembly functionality with three core features:

1. **Variable Extraction** - Extract variables from Go templates
2. **Default Value Extraction** - Extract variables with default values from custom functions
3. **Template Rendering** - Render templates with provided variables

## 🚀 Quick Start

### 1. Build the WebAssembly Module

```bash
# Build the WASM binary
GOOS=js GOARCH=wasm go build -o main.wasm
```


## 🔧 Supported Template Features

### Variables
- Simple fields: `{{.Name}}`
- Nested fields: `{{.User.Name}}`
- Custom functions: `{{getv "key" "default"}}`

### Custom Functions
- `getv "key" "default"` - Get value with default
- `exists "key"` - Check if key exists
- `get "key"` - Get value (error if not found)
- `jsonv "key"` - Get JSON value as object

### Control Structures
- Conditionals: `{{if .Condition}}...{{end}}`
- Loops: `{{range .Items}}...{{end}}`
- Context: `{{with .User}}...{{end}}`
- Else branches: `{{if .Active}}...{{else}}...{{end}}`

## 📁 Files Structure

```
├── main.wasm              # WebAssembly module (build this)
├── wasm_exec.js           # Go WASM support (copied from Go)
├── main.go                # WASM entry point
├── template_parser.go     # Core parsing logic
├── wasm_handlers.go       # JavaScript interface
└── README_DEMO.md         # This file
```

## 🛠️ Development

### Building WASM
```bash
# Build for WebAssembly
GOOS=js GOARCH=wasm go build -o main.wasm

# Build with optimization
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o main.wasm
```

### Testing
```bash
# Run unit tests
go test -v template_parser_test.go template_parser_pure_test.go template_parser.go

# Test with coverage
go test -cover -coverprofile=coverage.out template_parser_test.go template_parser_pure_test.go template_parser.go
```

## 🤝 Contributing

Feel free to:
- Add more template examples
- Improve the UI/UX
- Add new features
- Report bugs
- Suggest improvements