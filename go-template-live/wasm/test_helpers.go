//go:build !js
// +build !js

package main

import (
	"text/template/parse"
)

// TestHelper provides utilities for testing both build configurations
type TestHelper struct{}

// NewTestHelper creates a new test helper
func NewTestHelper() *TestHelper {
	return &TestHelper{}
}

// NewParserWithCustomFunctions creates a parser with custom functions enabled
// This simulates the custom build by manually registering custom functions
func (h *TestHelper) NewParserWithCustomFunctions() *Parser {
	registry := NewFunctionRegistry()

	// Manually register custom functions for testing
	// In the actual build, this is done by init() in functions_custom.go
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "getv",
		Description:           "Get variable value with optional default (Confd-style)",
		Handler:               func(key string, v ...string) string { return "" },
		Extractor:             extractGetvVariables,
		ExtractorWithDefaults: extractGetvVariablesWithDefaults,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "exists",
		Description:           "Check if variable exists (Confd-style)",
		Handler:               func(key string) bool { return false },
		Extractor:             extractKeyArgVariable,
		ExtractorWithDefaults: extractKeyArgVariableInfo,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "get",
		Description:           "Get variable value, returns error if not found (Confd-style)",
		Handler:               func(key string) (interface{}, error) { return nil, nil },
		Extractor:             extractKeyArgVariable,
		ExtractorWithDefaults: extractKeyArgVariableInfo,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "json",
		Description:           "Parse JSON variable and return as map (Confd-style)",
		Handler:               func(key string) (map[string]interface{}, error) { return nil, nil },
		Extractor:             extractKeyArgVariable,
		ExtractorWithDefaults: extractKeyArgVariableInfo,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "jsonArray",
		Description:           "Parse JSON variable and return as array (Confd-style)",
		Handler:               func(key string) ([]interface{}, error) { return nil, nil },
		Extractor:             extractKeyArgVariable,
		ExtractorWithDefaults: extractKeyArgVariableInfo,
	})

	return NewParser(registry)
}

// NewParserWithConfdFunctions creates a parser with Confd-style functions enabled
// This simulates the confd build by manually registering Confd functions
func (h *TestHelper) NewParserWithConfdFunctions() *Parser {
	registry := NewFunctionRegistry()

	// Manually register Confd functions for testing
	// In the actual build, this is done by init() in functions_confd.go

	// Pure utility functions (no variable extraction)
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "base",
		Description:           "Returns the last element of path",
		Handler:               func(s string) string { return "" },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "split",
		Description:           "Splits a string into substrings separated by separator",
		Handler:               func(s, sep string) []string { return []string{} },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "dir",
		Description:           "Returns all but the last element of path",
		Handler:               func(s string) string { return "" },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "map",
		Description:           "Create a map from key-value pairs",
		Handler:               func(values ...interface{}) (map[string]interface{}, error) { return nil, nil },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "join",
		Description:           "Joins array elements into a string with separator",
		Handler:               func(elems []string, sep string) string { return "" },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "datetime",
		Description:           "Returns current time",
		Handler:               func() interface{} { return nil },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "toUpper",
		Description:           "Converts string to uppercase",
		Handler:               func(s string) string { return "" },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "toLower",
		Description:           "Converts string to lowercase",
		Handler:               func(s string) string { return "" },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "replace",
		Description:           "Replaces old string with new string",
		Handler:               func(s, old, new string, n int) string { return "" },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "contains",
		Description:           "Checks if string contains substring",
		Handler:               func(s, substr string) bool { return false },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "base64Encode",
		Description:           "Base64 encodes a string",
		Handler:               func(data string) string { return "" },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "base64Decode",
		Description:           "Base64 decodes a string",
		Handler:               func(data string) (string, error) { return "", nil },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "trimSuffix",
		Description:           "Trims suffix from string",
		Handler:               func(s, suffix string) string { return "" },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "parseBool",
		Description:           "Parses string to boolean",
		Handler:               func(str string) (bool, error) { return false, nil },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "reverse",
		Description:           "Reverses an array",
		Handler:               func(values interface{}) interface{} { return nil },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "add",
		Description:           "Adds two numbers",
		Handler:               func(a, b int) int { return 0 },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "sub",
		Description:           "Subtracts two numbers",
		Handler:               func(a, b int) int { return 0 },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "div",
		Description:           "Divides two numbers",
		Handler:               func(a, b int) int { return 0 },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "mod",
		Description:           "Modulo operation on two numbers",
		Handler:               func(a, b int) int { return 0 },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "mul",
		Description:           "Multiplies two numbers",
		Handler:               func(a, b int) int { return 0 },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "seq",
		Description:           "Generates sequence of integers",
		Handler:               func(first, last int) []int { return []int{} },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "atoi",
		Description:           "Converts string to integer",
		Handler:               func(s string) (int, error) { return 0, nil },
		Extractor:             extractNoVariablesTest,
		ExtractorWithDefaults: extractNoVariablesInfoTest,
	})

	// JSON functions (these DO extract variables)
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "json",
		Description:           "Parse JSON variable and return as map",
		Handler:               func(key string) (map[string]interface{}, error) { return nil, nil },
		Extractor:             extractSingleStringArgVariable,
		ExtractorWithDefaults: extractSingleStringArgVariableInfo,
	})

	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "jsonArray",
		Description:           "Parse JSON variable and return as array",
		Handler:               func(key string) ([]interface{}, error) { return nil, nil },
		Extractor:             extractSingleStringArgVariable,
		ExtractorWithDefaults: extractSingleStringArgVariableInfo,
	})

	return NewParser(registry)
}

// NewParserWithOfficialFunctions creates a parser with only official functions
// This simulates the official build with no custom functions registered
func (h *TestHelper) NewParserWithOfficialFunctions() *Parser {
	registry := NewFunctionRegistry()
	// Don't register any custom functions
	return NewParser(registry)
}

// ShouldTestCustomFunctions returns whether custom functions should be tested
func (h *TestHelper) ShouldTestCustomFunctions() bool {
	return true
}

// Extractors for Confd functions (defined here since they're only used in tests)
// Using different names to avoid conflicts with functions_confd.go
func extractNoVariablesTest(args []parse.Node, cycle int) ([]string, error) {
	return nil, nil
}

func extractNoVariablesInfoTest(args []parse.Node, cycle int) ([]VariableInfo, error) {
	return nil, nil
}
