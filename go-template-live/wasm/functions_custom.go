//go:build custom
// +build custom

// This file contains the core implementations of custom functions
// Tag: custom (works for both js && custom WASM builds and !js && custom tests)

package main

import (
	"encoding/json"
	"fmt"
	"text/template"
	"text/template/parse"
)

// registerCustomFunctions registers all custom template functions
// This is called by both WASM (via init in main_custom.go) and tests
func registerCustomFunctions() {
	registry := GetGlobalRegistry()

	// getv - Get variable value with optional default
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "getv",
		Description:           "Get variable value with optional default (Confd-style)",
		Handler:               getvMinimalHandler,
		Extractor:             extractGetvVariables,
		ExtractorWithDefaults: extractGetvVariablesWithDefaults,
	})

	// exists - Check if variable exists (no default value support)
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "exists",
		Description:           "Check if variable exists (Confd-style)",
		Handler:               existsMinimalHandler,
		Extractor:             extractKeyArgVariable,
		ExtractorWithDefaults: extractKeyArgVariableInfo,
	})

	// get - Get variable value (errors if not found, no default value support)
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "get",
		Description:           "Get variable value, returns error if not found (Confd-style)",
		Handler:               getMinimalHandler,
		Extractor:             extractKeyArgVariable,
		ExtractorWithDefaults: extractKeyArgVariableInfo,
	})

	// json - Parse JSON variable (no default value support)
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "json",
		Description:           "Parse JSON variable and return as map (Confd-style)",
		Handler:               jsonMinimalHandler,
		Extractor:             extractKeyArgVariable,
		ExtractorWithDefaults: extractKeyArgVariableInfo,
	})

	// jsonArray - Parse JSON variable and return as array (no default value support)
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "jsonArray",
		Description:           "Parse JSON variable and return as array (Confd-style)",
		Handler:               jsonArrayMinimalHandler,
		Extractor:             extractKeyArgVariable,
		ExtractorWithDefaults: extractKeyArgVariableInfo,
	})
}

// Minimal handlers for parsing (don't need actual variable values)
func getvMinimalHandler(key string, v ...string) string {
	return ""
}

func existsMinimalHandler(key string) bool {
	return false
}

func getMinimalHandler(key string) (interface{}, error) {
	return nil, nil
}

func jsonMinimalHandler(key string) (map[string]interface{}, error) {
	return nil, nil
}

func jsonArrayMinimalHandler(key string) ([]interface{}, error) {
	return nil, nil
}

// Actual handlers for rendering (use variable values)
func getvRenderHandler(variables map[string]interface{}) func(key string, v ...string) string {
	return func(key string, v ...string) string {
		if val, exists := variables[key]; exists {
			if strVal, ok := val.(string); ok && strVal != "" {
				return strVal
			}
		}
		if len(v) > 0 {
			return v[0] // return default value
		}
		return ""
	}
}

func existsRenderHandler(variables map[string]interface{}) func(key string) bool {
	return func(key string) bool {
		_, exists := variables[key]
		return exists
	}
}

func getRenderHandler(variables map[string]interface{}) func(key string) (interface{}, error) {
	return func(key string) (interface{}, error) {
		if val, exists := variables[key]; exists {
			return val, nil
		}
		return nil, fmt.Errorf("key %s not found", key)
	}
}

func jsonRenderHandler(variables map[string]interface{}) func(key string) (map[string]interface{}, error) {
	return func(key string) (map[string]interface{}, error) {
		if val, exists := variables[key]; exists {
			if strVal, ok := val.(string); ok {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(strVal), &result)
				return result, err
			}
		}
		return nil, fmt.Errorf("key %s not found", key)
	}
}

func jsonArrayRenderHandler(variables map[string]interface{}) func(key string) ([]interface{}, error) {
	return func(key string) ([]interface{}, error) {
		if val, exists := variables[key]; exists {
			if strVal, ok := val.(string); ok {
				var result []interface{}
				err := json.Unmarshal([]byte(strVal), &result)
				return result, err
			}
		}
		return nil, fmt.Errorf("key %s not found", key)
	}
}

// Variable extractors
// Note: Only getv supports default values. Other functions (exists, get, json)
// just extract the key name without defaults.

// extractKeyArgVariable extracts a single key argument (used by all custom functions)
func extractKeyArgVariable(args []parse.Node, cycle int) ([]string, error) {
	return extractStringArgVariable(args, cycle, 1)
}

// extractKeyArgVariableInfo extracts key as VariableInfo without defaults
func extractKeyArgVariableInfo(args []parse.Node, cycle int) ([]VariableInfo, error) {
	return extractStringArgVariableWithDefaults(args, cycle, 1, -1)
}

// getv is special - it supports default values
func extractGetvVariables(args []parse.Node, cycle int) ([]string, error) {
	return extractStringArgVariable(args, cycle, 1)
}

func extractGetvVariablesWithDefaults(args []parse.Node, cycle int) ([]VariableInfo, error) {
	// getv supports default value as second argument (index 2)
	return extractStringArgVariableWithDefaults(args, cycle, 1, 2)
}

// GetCustomRenderFuncMap returns a function map with all custom functions for rendering
// This is used by both the WASM build (via CreateRenderFuncMap) and tests
func GetCustomRenderFuncMap(variables map[string]interface{}) template.FuncMap {
	return template.FuncMap{
		"getv":      getvRenderHandler(variables),
		"exists":    existsRenderHandler(variables),
		"get":       getRenderHandler(variables),
		"json":      jsonRenderHandler(variables),
		"jsonArray": jsonArrayRenderHandler(variables),
	}
}
