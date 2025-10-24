package main

import (
	"text/template/parse"
)

// extractStringArgVariable is a helper function to extract variable from string arguments
// This is used by many custom functions that take a key as their first argument
func extractStringArgVariable(args []parse.Node, cycle int, argIndex int) ([]string, error) {
	var result []string
	if len(args) > argIndex {
		item := args[argIndex]
		if item.Type() == parse.NodeString {
			stringNode := item.(*parse.StringNode)
			result = append(result, stringNode.Text)
		} else {
			// For complex expressions, we need to recursively extract variables
			parser := NewParser(globalRegistry)
			sonResult, err := parser.getFieldFromNode(item, cycle)
			if err != nil {
				return nil, err
			}
			result = append(result, sonResult...)
		}
	}
	return result, nil
}

// extractStringArgVariableWithDefaults extracts variables with default values
// argIndex: the index of the variable name argument (usually 1, after function name)
// defaultArgIndex: the index of the default value argument (usually 2)
func extractStringArgVariableWithDefaults(args []parse.Node, cycle int, argIndex int, defaultArgIndex int) ([]VariableInfo, error) {
	var result []VariableInfo
	if len(args) > argIndex {
		item := args[argIndex]
		if item.Type() == parse.NodeString {
			stringNode := item.(*parse.StringNode)
			varInfo := VariableInfo{
				Name: stringNode.Text,
			}
			// Check for default value
			if defaultArgIndex > 0 && len(args) > defaultArgIndex && args[defaultArgIndex].Type() == parse.NodeString {
				defaultValueNode := args[defaultArgIndex].(*parse.StringNode)
				varInfo.DefaultValue = defaultValueNode.Text
			}
			result = append(result, varInfo)
		} else {
			// For complex expressions, extract variable names without defaults
			parser := NewParser(globalRegistry)
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

// Global registry instance
// This will be populated by init() functions in function implementation files
var globalRegistry = NewFunctionRegistry()

// GetGlobalRegistry returns the global function registry
func GetGlobalRegistry() *FunctionRegistry {
	return globalRegistry
}
