//go:build !js
// +build !js

package main

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
