//go:build js
// +build js

package main

import (
	"encoding/json"
	"strings"
	"syscall/js"
	"text/template"
)

// WASMHandler handles WASM/JavaScript interface operations
type WASMHandler struct {
	parser *Parser
}

// NewWASMHandler creates a new WASM handler
func NewWASMHandler() *WASMHandler {
	return NewWASMHandlerWithConfig(DefaultBuildConfig())
}

// NewWASMHandlerWithConfig creates a new WASM handler with specific configuration
func NewWASMHandlerWithConfig(config *BuildConfig) *WASMHandler {
	functionMatcher := &DefaultFunctionMatcher{}
	parser := NewParserWithConfig(functionMatcher, config)
	return &WASMHandler{
		parser: parser,
	}
}

// ExtractVariables extracts variables with default values - main function exposed to JavaScript
func (h *WASMHandler) ExtractVariables(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return jsError("Missing template content parameter")
	}

	templateContent := args[0].String()
	fileName := "template.tmpl"
	if len(args) > 1 {
		fileName = args[1].String()
	}

	variables, err := h.parser.ExtractVariablesWithDefaults(fileName, templateContent)
	if err != nil {
		return jsError("Failed to extract variables: " + err.Error())
	}

	jsonData, err := json.Marshal(variables)
	if err != nil {
		return jsError("Failed to marshal variables to JSON: " + err.Error())
	}

	return js.ValueOf(string(jsonData))
}

// ExtractVariablesSimple returns only variable names (without defaults)
func (h *WASMHandler) ExtractVariablesSimple(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return jsError("Missing template content parameter")
	}

	templateContent := args[0].String()
	fileName := "template.tmpl"
	if len(args) > 1 {
		fileName = args[1].String()
	}

	variables, err := h.parser.ExtractVariables(fileName, templateContent)
	if err != nil {
		return jsError("Failed to extract variables: " + err.Error())
	}

	jsonData, err := json.Marshal(variables)
	if err != nil {
		return jsError("Failed to marshal variables to JSON: " + err.Error())
	}

	return js.ValueOf(string(jsonData))
}

// RenderTemplate renders a template with provided variable values
func (h *WASMHandler) RenderTemplate(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return jsError("Missing template content or variables parameter")
	}

	templateContent := args[0].String()
	variablesJSON := args[1].String()

	var variables map[string]interface{}
	err := json.Unmarshal([]byte(variablesJSON), &variables)
	if err != nil {
		return jsError("Failed to parse variables JSON: " + err.Error())
	}

	functionHandler := NewFunctionHandler(variables)
	funcs := functionHandler.CreateFuncMapWithConfig(h.parser.functionRegistry.config)

	tmpl, err := template.New("template").Funcs(funcs).Parse(templateContent)
	if err != nil {
		return jsError("Failed to parse template: " + err.Error())
	}

	var result strings.Builder
	err = tmpl.Execute(&result, variables)
	if err != nil {
		return jsError("Failed to execute template: " + err.Error())
	}

	return js.ValueOf(result.String())
}

// RegisterCallbacks registers the Go functions to be called from JavaScript
func (h *WASMHandler) RegisterCallbacks() {
	js.Global().Set("extractTemplateVariables", js.FuncOf(h.ExtractVariables))
	js.Global().Set("extractTemplateVariablesSimple", js.FuncOf(h.ExtractVariablesSimple))
	js.Global().Set("renderTemplateWithValues", js.FuncOf(h.RenderTemplate))
}

// jsError creates a JavaScript error object
func jsError(message string) map[string]interface{} {
	return map[string]interface{}{
		"error": message,
	}
}
