package main

import (
	"text/template/parse"
)

// extractArgVariable is a helper function to extract variables from function arguments
// argIndex: the index of the argument to extract from
// treatStringLiteralsAsVarNames: if true, string literals like "myvar" are treated as variable names (for json, getv);
//
//	if false, only field accesses like .myvar are extracted (for base, toUpper, etc.)
func extractArgVariable(args []parse.Node, cycle int, argIndex int, treatStringLiteralsAsVarNames bool) ([]string, error) {
	var result []string
	if len(args) > argIndex {
		item := args[argIndex]
		if item.Type() == parse.NodeString {
			if treatStringLiteralsAsVarNames {
				// String literals are variable names to look up (e.g., json "myvar")
				stringNode := item.(*parse.StringNode)
				result = append(result, stringNode.Text)
			}
			// else: string literals are just data, skip them
		} else {
			// For complex expressions, recursively extract variables
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

// extractArgVariableWithDefaults extracts variables with default values
// argIndex: the index of the variable name argument (usually 1, after function name)
// defaultArgIndex: the index of the default value argument (usually 2, -1 for no default)
// treatStringLiteralsAsVarNames: if true, string literals are treated as variable names
func extractArgVariableWithDefaults(args []parse.Node, cycle int, argIndex int, defaultArgIndex int, treatStringLiteralsAsVarNames bool) ([]VariableInfo, error) {
	var result []VariableInfo
	if len(args) > argIndex {
		item := args[argIndex]
		if item.Type() == parse.NodeString {
			if treatStringLiteralsAsVarNames {
				// String literals are variable names to look up
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
			}
			// else: string literals are just data, skip them
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

// Legacy wrappers for backward compatibility

// extractStringArgVariable treats string literals as variable names (for json, getv, etc.)
func extractStringArgVariable(args []parse.Node, cycle int, argIndex int) ([]string, error) {
	return extractArgVariable(args, cycle, argIndex, true)
}

func extractStringArgVariableWithDefaults(args []parse.Node, cycle int, argIndex int, defaultArgIndex int) ([]VariableInfo, error) {
	return extractArgVariableWithDefaults(args, cycle, argIndex, defaultArgIndex, true)
}

// Global registry instance
// This will be populated by init() functions in function implementation files
var globalRegistry = NewFunctionRegistry()

// GetGlobalRegistry returns the global function registry
func GetGlobalRegistry() *FunctionRegistry {
	return globalRegistry
}
