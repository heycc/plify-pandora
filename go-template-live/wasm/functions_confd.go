//go:build confd
// +build confd

// This file contains the core implementations of Confd-style functions
// Tag: confd (works for both js && confd WASM builds and !js && confd tests)

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"strings"
	"text/template"
	"text/template/parse"
	"time"
)

// registerConfdFunctions registers all Confd-style template functions
// This is called by both WASM (via init in functions_confd.go) and tests
func registerConfdFunctions() {
	registry := GetGlobalRegistry()

	// Custom functions (getv, exists, get)
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

	// base - Base function (path.Base) - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "base",
		Description:           "Returns the last element of path",
		Handler:               baseMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// split - Split function (strings.Split) - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "split",
		Description:           "Splits a string into substrings separated by separator",
		Handler:               splitMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// json - Parse JSON object
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "json",
		Description:           "Parse JSON variable and return as map",
		Handler:               jsonConfdMinimalHandler,
		Extractor:             extractSingleStringArgVariable,
		ExtractorWithDefaults: extractSingleStringArgVariableInfo,
	})

	// jsonArray - Parse JSON array
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "jsonArray",
		Description:           "Parse JSON variable and return as array",
		Handler:               jsonArrayConfdMinimalHandler,
		Extractor:             extractSingleStringArgVariable,
		ExtractorWithDefaults: extractSingleStringArgVariableInfo,
	})

	// dir - Directory function (path.Dir) - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "dir",
		Description:           "Returns all but the last element of path",
		Handler:               dirMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// map - Create map from key-value pairs - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "map",
		Description:           "Create a map from key-value pairs",
		Handler:               mapMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// join - Join function (strings.Join) - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "join",
		Description:           "Joins array elements into a string with separator",
		Handler:               joinMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// datetime - Current time - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "datetime",
		Description:           "Returns current time",
		Handler:               datetimeMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// toUpper - Convert to uppercase - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "toUpper",
		Description:           "Converts string to uppercase",
		Handler:               toUpperMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// toLower - Convert to lowercase - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "toLower",
		Description:           "Converts string to lowercase",
		Handler:               toLowerMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// replace - Replace substrings - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "replace",
		Description:           "Replaces old string with new string",
		Handler:               replaceMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// contains - Check if string contains substring - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "contains",
		Description:           "Checks if string contains substring",
		Handler:               containsMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// base64Encode - Base64 encode - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "base64Encode",
		Description:           "Base64 encodes a string",
		Handler:               base64EncodeMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// base64Decode - Base64 decode - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "base64Decode",
		Description:           "Base64 decodes a string",
		Handler:               base64DecodeMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// trimSuffix - Trim suffix - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "trimSuffix",
		Description:           "Trims suffix from string",
		Handler:               trimSuffixMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// parseBool - Parse boolean - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "parseBool",
		Description:           "Parses string to boolean",
		Handler:               parseBoolMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// reverse - Reverse array - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "reverse",
		Description:           "Reverses an array",
		Handler:               reverseMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// add - Add numbers - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "add",
		Description:           "Adds two numbers",
		Handler:               addMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// sub - Subtract numbers - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "sub",
		Description:           "Subtracts two numbers",
		Handler:               subMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// div - Divide numbers - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "div",
		Description:           "Divides two numbers",
		Handler:               divMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// mod - Modulo operation - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "mod",
		Description:           "Modulo operation on two numbers",
		Handler:               modMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// mul - Multiply numbers - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "mul",
		Description:           "Multiplies two numbers",
		Handler:               mulMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// seq - Generate sequence - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "seq",
		Description:           "Generates sequence of integers",
		Handler:               seqMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// atoi - String to integer - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "atoi",
		Description:           "Converts string to integer",
		Handler:               atoiMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})
}

// Minimal handlers for parsing (don't need actual variable values)
func baseMinimalHandler(s string) string                                      { return "" }
func splitMinimalHandler(s, sep string) []string                              { return []string{} }
func jsonConfdMinimalHandler(data string) (map[string]interface{}, error)     { return nil, nil }
func jsonArrayConfdMinimalHandler(data string) ([]interface{}, error)         { return nil, nil }
func dirMinimalHandler(s string) string                                       { return "" }
func mapMinimalHandler(values ...interface{}) (map[string]interface{}, error) { return nil, nil }
func joinMinimalHandler(elems []string, sep string) string                    { return "" }
func datetimeMinimalHandler() time.Time                                       { return time.Time{} }
func toUpperMinimalHandler(s string) string                                   { return "" }
func toLowerMinimalHandler(s string) string                                   { return "" }
func replaceMinimalHandler(s, old, new string, n int) string                  { return "" }
func containsMinimalHandler(s, substr string) bool                            { return false }
func base64EncodeMinimalHandler(data string) string                           { return "" }
func base64DecodeMinimalHandler(data string) (string, error)                  { return "", nil }
func trimSuffixMinimalHandler(s, suffix string) string                        { return "" }
func parseBoolMinimalHandler(str string) (bool, error)                        { return false, nil }
func reverseMinimalHandler(values interface{}) interface{}                    { return nil }
func addMinimalHandler(a, b int) int                                          { return 0 }
func subMinimalHandler(a, b int) int                                          { return 0 }
func divMinimalHandler(a, b int) int                                          { return 0 }
func modMinimalHandler(a, b int) int                                          { return 0 }
func mulMinimalHandler(a, b int) int                                          { return 0 }
func seqMinimalHandler(first, last int) []int                                 { return []int{} }
func atoiMinimalHandler(s string) (int, error)                                { return 0, nil }

// Variable extractors (used during template parsing)
func extractNoVariables(args []parse.Node, cycle int) ([]string, error) {
	return []string{}, nil
}

func extractNoVariablesInfo(args []parse.Node, cycle int) ([]VariableInfo, error) {
	return []VariableInfo{}, nil
}

// Only json and jsonArray need variable extraction since they access variables
func extractSingleStringArgVariable(args []parse.Node, cycle int) ([]string, error) {
	return extractStringArgVariable(args, cycle, 1)
}

func extractSingleStringArgVariableInfo(args []parse.Node, cycle int) ([]VariableInfo, error) {
	return extractStringArgVariableWithDefaults(args, cycle, 1, -1)
}

// Extract variables from first argument (for transformation functions like base, toUpper, etc.)
// These functions operate on their first argument, which may be a variable reference
// Unlike extractStringArgVariable, this does NOT treat string literals as variable names
func extractFirstArgVariable(args []parse.Node, cycle int) ([]string, error) {
	return extractArgVariable(args, cycle, 1, false)
}

func extractFirstArgVariableInfo(args []parse.Node, cycle int) ([]VariableInfo, error) {
	return extractArgVariableWithDefaults(args, cycle, 1, -1, false)
}

// GetConfdRenderFuncMap returns a function map with all Confd-style functions for rendering
// This is used by both the WASM build (via CreateRenderFuncMap) and tests
func GetConfdRenderFuncMap(variables map[string]interface{}) template.FuncMap {
	return template.FuncMap{
		// Custom functions (getv, exists, get)
		"getv": func(key string, v ...string) string {
			if val, exists := variables[key]; exists {
				if strVal, ok := val.(string); ok && strVal != "" {
					return strVal
				}
			}
			if len(v) > 0 {
				return v[0] // return default value
			}
			return ""
		},
		"exists": func(key string) bool {
			_, exists := variables[key]
			return exists
		},
		"get": func(key string) (interface{}, error) {
			if val, exists := variables[key]; exists {
				return val, nil
			}
			return nil, fmt.Errorf("key %s not found", key)
		},
		// Confd functions
		"base":         func(s string) string { return path.Base(s) },
		"split":        func(s, sep string) []string { return strings.Split(s, sep) },
		"dir":          func(s string) string { return path.Dir(s) },
		"join":         func(elems []string, sep string) string { return strings.Join(elems, sep) },
		"datetime":     func() time.Time { return time.Now() },
		"toUpper":      func(s string) string { return strings.ToUpper(s) },
		"toLower":      func(s string) string { return strings.ToLower(s) },
		"replace":      func(s, old, new string, n int) string { return strings.Replace(s, old, new, n) },
		"contains":     func(s, substr string) bool { return strings.Contains(s, substr) },
		"base64Encode": func(data string) string { return base64.StdEncoding.EncodeToString([]byte(data)) },
		"base64Decode": func(data string) (string, error) {
			s, err := base64.StdEncoding.DecodeString(data)
			return string(s), err
		},
		"trimSuffix": func(s, suffix string) string { return strings.TrimSuffix(s, suffix) },
		"parseBool":  func(str string) (bool, error) { return strconv.ParseBool(str) },
		"add":        func(a, b int) int { return a + b },
		"sub":        func(a, b int) int { return a - b },
		"div":        func(a, b int) int { return a / b },
		"mod":        func(a, b int) int { return a % b },
		"mul":        func(a, b int) int { return a * b },
		"seq": func(first, last int) []int {
			var result []int
			for i := first; i <= last; i++ {
				result = append(result, i)
			}
			return result
		},
		"atoi": func(s string) (int, error) { return strconv.Atoi(s) },
		"map": func(values ...interface{}) (map[string]interface{}, error) {
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, nil
				}
				if i+1 < len(values) {
					dict[key] = values[i+1]
				}
			}
			return dict, nil
		},
		"reverse": func(values interface{}) interface{} {
			switch v := values.(type) {
			case []string:
				for left, right := 0, len(v)-1; left < right; left, right = left+1, right-1 {
					v[left], v[right] = v[right], v[left]
				}
				return v
			case []interface{}:
				for left, right := 0, len(v)-1; left < right; left, right = left+1, right-1 {
					v[left], v[right] = v[right], v[left]
				}
				return v
			}
			return values
		},
		// JSON functions that look up variables
		"json": func(key string) (map[string]interface{}, error) {
			if val, ok := variables[key]; ok {
				if str, ok := val.(string); ok {
					var result map[string]interface{}
					err := json.Unmarshal([]byte(str), &result)
					return result, err
				}
			}
			return nil, nil
		},
		"jsonArray": func(key string) ([]interface{}, error) {
			if val, ok := variables[key]; ok {
				if str, ok := val.(string); ok {
					var result []interface{}
					err := json.Unmarshal([]byte(str), &result)
					return result, err
				}
			}
			return nil, nil
		},
	}
}

// Minimal handlers for custom functions
func getvMinimalHandler(key string, v ...string) string {
	return ""
}

func existsMinimalHandler(key string) bool {
	return false
}

func getMinimalHandler(key string) (interface{}, error) {
	return nil, nil
}

// Variable extractors for custom functions
// extractKeyArgVariable extracts a single key argument (used by custom functions)
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
