//go:build js && confd
// +build js,confd

// This file contains WASM-specific wiring for Confd functions
// The actual implementations are in functions_confd_core.go

package main

func init() {
	// Register Confd-style functions on initialization
	// This only happens when building WASM with the "js && confd" tags
	// The registerConfdFunctions() function is in functions_confd_core.go
	registerConfdFunctions()
}

// CreateRenderFuncMap creates function map with actual variable values for rendering Confd functions
// This delegates to GetConfdRenderFuncMap from functions_confd_core.go
func CreateRenderFuncMap(variables map[string]interface{}) map[string]interface{} {
	funcMap := GetConfdRenderFuncMap(variables)
	// Convert template.FuncMap to map[string]interface{}
	result := make(map[string]interface{}, len(funcMap))
	for k, v := range funcMap {
		result[k] = v
	}
	return result
}
