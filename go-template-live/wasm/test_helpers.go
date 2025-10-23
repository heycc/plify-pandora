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
func (h *TestHelper) NewParserWithCustomFunctions() *Parser {
	matcher := NewDefaultFunctionMatcher()
	return NewParserWithConfig(matcher, DefaultBuildConfig())
}

// NewParserWithOfficialFunctions creates a parser with only official functions
func (h *TestHelper) NewParserWithOfficialFunctions() *Parser {
	matcher := NewOfficialFunctionMatcher()
	return NewParserWithConfig(matcher, OfficialBuildConfig())
}

// ShouldTestCustomFunctions returns whether custom functions should be tested
// This helps tests skip custom function tests when running in official mode
func (h *TestHelper) ShouldTestCustomFunctions() bool {
	// For now, always test custom functions in unit tests
	// In the future, we could use build tags to control this
	return true
}