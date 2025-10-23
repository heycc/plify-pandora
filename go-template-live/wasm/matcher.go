package main

// FunctionMatcher defines the interface for matching custom template functions
// This abstraction allows different matching strategies, similar to how Confd
// supports different backends for configuration management
type FunctionMatcher interface {
	MatchCustomFunc(funcName string) bool
	GetSupportedFunctions() []string
}

// DefaultFunctionMatcher implements the default function matching logic
// This follows the Confd pattern of having a standard set of template functions
type DefaultFunctionMatcher struct {
	supportedFunctions map[string]bool
}

// NewDefaultFunctionMatcher creates a new default function matcher
func NewDefaultFunctionMatcher() *DefaultFunctionMatcher {
	return &DefaultFunctionMatcher{
		supportedFunctions: map[string]bool{
			"exists": true,
			"get":    true,
			"getv":   true,
			"jsonv":  true,
		},
	}
}

// NewOfficialFunctionMatcher creates a function matcher for official mode
// This matcher doesn't recognize any custom functions
func NewOfficialFunctionMatcher() *DefaultFunctionMatcher {
	return &DefaultFunctionMatcher{
		supportedFunctions: map[string]bool{},
	}
}

// MatchCustomFunc checks if a function should be recognized as a parameter-getting method
func (m *DefaultFunctionMatcher) MatchCustomFunc(funcName string) bool {
	return m.supportedFunctions[funcName]
}

// GetSupportedFunctions returns the list of supported function names
func (m *DefaultFunctionMatcher) GetSupportedFunctions() []string {
	functions := make([]string, 0, len(m.supportedFunctions))
	for funcName := range m.supportedFunctions {
		functions = append(functions, funcName)
	}
	return functions
}

// RegistryFunctionMatcher uses a function registry for matching
// This provides more flexibility for adding new functions dynamically
type RegistryFunctionMatcher struct {
	registry *FunctionRegistry
}

// NewRegistryFunctionMatcher creates a new registry-based function matcher
func NewRegistryFunctionMatcher(registry *FunctionRegistry) *RegistryFunctionMatcher {
	return &RegistryFunctionMatcher{
		registry: registry,
	}
}

// MatchCustomFunc checks if a function exists in the registry
func (m *RegistryFunctionMatcher) MatchCustomFunc(funcName string) bool {
	_, exists := m.registry.GetFunction(funcName)
	return exists
}

// GetSupportedFunctions returns the list of supported function names from the registry
func (m *RegistryFunctionMatcher) GetSupportedFunctions() []string {
	functions := make([]string, 0)
	// Note: We can't directly access the registry's internal map here
	// This would need to be implemented differently if we want to use this matcher
	return functions
}

// CompositeFunctionMatcher allows combining multiple matchers
// This is useful for supporting both built-in and custom functions
// Similar to how Confd supports multiple template function sources
type CompositeFunctionMatcher struct {
	matchers []FunctionMatcher
}

// NewCompositeFunctionMatcher creates a new composite function matcher
func NewCompositeFunctionMatcher(matchers ...FunctionMatcher) *CompositeFunctionMatcher {
	return &CompositeFunctionMatcher{
		matchers: matchers,
	}
}

// MatchCustomFunc checks if any of the matchers recognize the function
func (m *CompositeFunctionMatcher) MatchCustomFunc(funcName string) bool {
	for _, matcher := range m.matchers {
		if matcher.MatchCustomFunc(funcName) {
			return true
		}
	}
	return false
}

// GetSupportedFunctions returns the combined list of supported function names
func (m *CompositeFunctionMatcher) GetSupportedFunctions() []string {
	allFunctions := make(map[string]bool)
	for _, matcher := range m.matchers {
		for _, funcName := range matcher.GetSupportedFunctions() {
			allFunctions[funcName] = true
		}
	}

	functions := make([]string, 0, len(allFunctions))
	for funcName := range allFunctions {
		functions = append(functions, funcName)
	}
	return functions
}
