package main

import (
	"encoding/json"
	"fmt"
	"text/template"
	"text/template/parse"
)

// FunctionDefinition defines a custom template function
// Inspired by Confd's approach to template functions
// https://github.com/kelseyhightower/confd

type FunctionDefinition struct {
	Name        string
	Description string
	Handler     interface{}
	Extractor   VariableExtractor
}

// VariableExtractor extracts variable names from function arguments
// This allows us to identify which variables are being used in templates
// similar to how Confd tracks template dependencies
type VariableExtractor func(args []parse.Node, cycle int) ([]string, error)

// VariableExtractorWithDefaults extracts variables with default values
type VariableExtractorWithDefaults func(args []parse.Node, cycle int) ([]VariableInfo, error)

// FunctionRegistry manages all custom template functions
type FunctionRegistry struct {
	functions map[string]*FunctionDefinition
}

// NewFunctionRegistry creates a new function registry
func NewFunctionRegistry() *FunctionRegistry {
	return &FunctionRegistry{
		functions: make(map[string]*FunctionDefinition),
	}
}

// RegisterFunction registers a new custom function
func (r *FunctionRegistry) RegisterFunction(def *FunctionDefinition) {
	r.functions[def.Name] = def
}

// GetFunction returns a function definition by name
func (r *FunctionRegistry) GetFunction(name string) (*FunctionDefinition, bool) {
	def, exists := r.functions[name]
	return def, exists
}

// GetFuncMap creates a template.FuncMap for rendering
func (r *FunctionRegistry) GetFuncMap(variables map[string]interface{}) template.FuncMap {
	funcMap := template.FuncMap{}
	for name, def := range r.functions {
		funcMap[name] = def.Handler
	}
	return funcMap
}

// GetMinimalFuncMap creates a minimal function map for parsing
func (r *FunctionRegistry) GetMinimalFuncMap() template.FuncMap {
	funcMap := template.FuncMap{}
	for name, def := range r.functions {
		// For parsing, we use minimal implementations that don't require actual variables
		funcMap[name] = def.Handler
	}
	return funcMap
}


// getvFunc implements getv function with default value support
// Similar to Confd's getv function
func getvFunc(variables map[string]interface{}) func(key string, v ...string) string {
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

// existsFunc implements exists function
// Similar to Confd's exists function
func existsFunc(variables map[string]interface{}) func(key string) bool {
	return func(key string) bool {
		_, exists := variables[key]
		return exists
	}
}

// getFunc implements get function
// Similar to Confd's get function
func getFunc(variables map[string]interface{}) func(key string) (interface{}, error) {
	return func(key string) (interface{}, error) {
		if val, exists := variables[key]; exists {
			return val, nil
		}
		return nil, fmt.Errorf("key %s not found", key)
	}
}

// jsonvFunc implements jsonv function
// Similar to Confd's JSON parsing capabilities
func jsonvFunc(variables map[string]interface{}) func(key string) (map[string]interface{}, error) {
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

// extractGetvVariables extracts variable names from getv function calls
func extractGetvVariables(args []parse.Node, cycle int) ([]string, error) {
	var result []string
	if len(args) > 1 {
		item := args[1]
		if item.Type() == parse.NodeString {
			stringNode := item.(*parse.StringNode)
			result = append(result, stringNode.Text)
		} else {
			// For complex expressions, we need to recursively extract variables
			// This is similar to how Confd handles nested template expressions
			parser := NewParser(NewDefaultFunctionMatcher())
			sonResult, err := parser.getFieldFromNode(item, cycle)
			if err != nil {
				return nil, err
			}
			result = append(result, sonResult...)
		}
	}
	return result, nil
}

// extractGetvVariablesWithDefaults extracts variables with default values from getv function calls
func extractGetvVariablesWithDefaults(args []parse.Node, cycle int) ([]VariableInfo, error) {
	var result []VariableInfo
	if len(args) > 1 {
		item := args[1]
		if item.Type() == parse.NodeString {
			stringNode := item.(*parse.StringNode)
			varInfo := VariableInfo{
				Name: stringNode.Text,
			}
			// Check for default value (third argument)
			if len(args) > 2 && args[2].Type() == parse.NodeString {
				defaultValueNode := args[2].(*parse.StringNode)
				varInfo.DefaultValue = defaultValueNode.Text
			}
			result = append(result, varInfo)
		} else {
			// For complex expressions, extract variable names without defaults
			parser := NewParser(NewDefaultFunctionMatcher())
			sonResult, err := parser.getFieldFromNode(item, cycle)
			if err != nil {
				return nil, err
			}
			for _, name := range sonResult {
				result = append(result, VariableInfo{Name: name})
			}
		}
	}
	return result, nil
}

// extractExistsVariables extracts variable names from exists function calls
func extractExistsVariables(args []parse.Node, cycle int) ([]string, error) {
	return extractGetvVariables(args, cycle)
}

// extractExistsVariablesWithDefaults extracts variables from exists function calls
func extractExistsVariablesWithDefaults(args []parse.Node, cycle int) ([]VariableInfo, error) {
	var result []VariableInfo
	if len(args) > 1 {
		item := args[1]
		if item.Type() == parse.NodeString {
			stringNode := item.(*parse.StringNode)
			result = append(result, VariableInfo{Name: stringNode.Text})
		} else {
			parser := NewParser(NewDefaultFunctionMatcher())
			sonResult, err := parser.getFieldFromNode(item, cycle)
			if err != nil {
				return nil, err
			}
			for _, name := range sonResult {
				result = append(result, VariableInfo{Name: name})
			}
		}
	}
	return result, nil
}

// extractGetVariables extracts variable names from get function calls
func extractGetVariables(args []parse.Node, cycle int) ([]string, error) {
	return extractGetvVariables(args, cycle)
}

// extractGetVariablesWithDefaults extracts variables from get function calls
func extractGetVariablesWithDefaults(args []parse.Node, cycle int) ([]VariableInfo, error) {
	return extractExistsVariablesWithDefaults(args, cycle)
}

// extractJsonvVariables extracts variable names from jsonv function calls
func extractJsonvVariables(args []parse.Node, cycle int) ([]string, error) {
	return extractGetvVariables(args, cycle)
}

// extractJsonvVariablesWithDefaults extracts variables from jsonv function calls
func extractJsonvVariablesWithDefaults(args []parse.Node, cycle int) ([]VariableInfo, error) {
	return extractExistsVariablesWithDefaults(args, cycle)
}

// DefaultFunctions returns the default function registry with all built-in functions
// This follows the Confd pattern of providing a standard set of template functions
func DefaultFunctions() *FunctionRegistry {
	registry := NewFunctionRegistry()

	// Register getv function - similar to Confd's getv
	registry.RegisterFunction(&FunctionDefinition{
		Name:        "getv",
		Description: "Get variable value with optional default",
		Handler:     func(key string, v ...string) string { return "" }, // Minimal implementation for parsing
		Extractor:   extractGetvVariables,
	})

	// Register exists function - similar to Confd's exists
	registry.RegisterFunction(&FunctionDefinition{
		Name:        "exists",
		Description: "Check if variable exists",
		Handler:     func(key string) bool { return false }, // Minimal implementation for parsing
		Extractor:   extractExistsVariables,
	})

	// Register get function - similar to Confd's get
	registry.RegisterFunction(&FunctionDefinition{
		Name:        "get",
		Description: "Get variable value (returns error if not found)",
		Handler:     func(key string) (interface{}, error) { return nil, nil }, // Minimal implementation for parsing
		Extractor:   extractGetVariables,
	})

	// Register jsonv function - similar to Confd's JSON capabilities
	registry.RegisterFunction(&FunctionDefinition{
		Name:        "jsonv",
		Description: "Parse JSON variable and return as map",
		Handler:     func(key string) (map[string]interface{}, error) { return nil, nil }, // Minimal implementation for parsing
		Extractor:   extractJsonvVariables,
	})

	return registry
}
