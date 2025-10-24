//go:build official
// +build official

package main

// This file is only included when building with the "official" tag
// Build with: GOOS=js GOARCH=wasm go build -tags official -o official.wasm .

// No custom functions are registered in official mode
// The init() function in functions_custom.go won't be called

// CreateRenderFuncMap creates an empty function map for official mode
// Only standard Go template functions will be available
func CreateRenderFuncMap(variables map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{}
}
