//go:build js && custom
// +build js,custom

// This file contains WASM-specific wiring for custom functions
// The actual implementations are in functions_custom_core.go

package main

func init() {
	// Register custom functions on initialization
	// This only happens when building WASM with the "js && custom" tags
	// The registerCustomFunctions() function is in functions_custom_core.go
	registerCustomFunctions()
}

// CreateRenderFuncMap creates function map with actual variable values for rendering custom functions
// This delegates to GetCustomRenderFuncMap from functions_custom_core.go
func CreateRenderFuncMap(variables map[string]interface{}) map[string]interface{} {
	funcMap := GetCustomRenderFuncMap(variables)
	// Convert template.FuncMap to map[string]interface{}
	result := make(map[string]interface{}, len(funcMap))
	for k, v := range funcMap {
		result[k] = v
	}
	return result
}
