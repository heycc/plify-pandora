package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"text/template/parse"
)

const maxDepth = 40

// VariableInfo stores variable information including name and default value
type VariableInfo struct {
	Name         string `json:"name"`
	DefaultValue string `json:"defaultValue"`
}

// FunctionMatcher handles custom function matching
type FunctionMatcher interface {
	MatchCustomFunc(funcName string) bool
}

// DefaultFunctionMatcher implements default function matching logic
type DefaultFunctionMatcher struct{}

// MatchCustomFunc checks if a function should be recognized as a parameter-getting method
func (m *DefaultFunctionMatcher) MatchCustomFunc(funcName string) bool {
	switch funcName {
	case "exists", "get", "getv", "jsonv":
		return true
	default:
		return false
	}
}

// FunctionHandler handles template function implementations
type FunctionHandler struct {
	variables map[string]interface{}
}

// NewFunctionHandler creates a new function handler
func NewFunctionHandler(variables map[string]interface{}) *FunctionHandler {
	return &FunctionHandler{
		variables: variables,
	}
}

// CreateFuncMap creates template function map for rendering
func (h *FunctionHandler) CreateFuncMap() template.FuncMap {
	return template.FuncMap{
		"getv": h.getvFunc,
		"exists": h.existsFunc,
		"get": h.getFunc,
		"jsonv": h.jsonvFunc,
	}
}

// getvFunc implements getv function with default value support
func (h *FunctionHandler) getvFunc(key string, v ...string) string {
	if val, exists := h.variables[key]; exists {
		if strVal, ok := val.(string); ok && strVal != "" {
			return strVal
		}
	}
	if len(v) > 0 {
		return v[0] // return default value
	}
	return ""
}

// existsFunc implements exists function
func (h *FunctionHandler) existsFunc(key string) bool {
	_, exists := h.variables[key]
	return exists
}

// getFunc implements get function
func (h *FunctionHandler) getFunc(key string) (interface{}, error) {
	if val, exists := h.variables[key]; exists {
		return val, nil
	}
	return nil, fmt.Errorf("key %s not found", key)
}

// jsonvFunc implements jsonv function
func (h *FunctionHandler) jsonvFunc(key string) (map[string]interface{}, error) {
	if val, exists := h.variables[key]; exists {
		if strVal, ok := val.(string); ok {
			var result map[string]interface{}
			err := json.Unmarshal([]byte(strVal), &result)
			return result, err
		}
	}
	return nil, fmt.Errorf("key %s not found", key)
}

// Parser handles template parsing and variable extraction
type Parser struct {
	functionMatcher FunctionMatcher
}

// NewParser creates a new template parser
func NewParser(matcher FunctionMatcher) *Parser {
	return &Parser{
		functionMatcher: matcher,
	}
}

// ExtractVariables extracts variable names from template content
func (p *Parser) ExtractVariables(fileName, fileContent string) ([]string, error) {
	funcs := p.createMinimalFuncMap()
	tmpl, err := template.New(fileName).Option("missingkey=error").Funcs(funcs).Parse(fileContent)
	if err != nil {
		return nil, fmt.Errorf("解析模板文件 %s 存在异常: %v", fileName, err)
	}

	result, err := p.getFieldFromNode(tmpl.Tree.Root, 0)
	if err != nil {
		return nil, fmt.Errorf("解析模板文件 %s 存在异常: %v", fileName, err)
	}

	return result, nil
}

// ExtractVariablesWithDefaults extracts variables with default values from template content
func (p *Parser) ExtractVariablesWithDefaults(fileName, fileContent string) ([]VariableInfo, error) {
	funcs := p.createMinimalFuncMap()
	tmpl, err := template.New(fileName).Option("missingkey=error").Funcs(funcs).Parse(fileContent)
	if err != nil {
		return nil, fmt.Errorf("解析模板文件 %s 存在异常: %v", fileName, err)
	}

	result, err := p.getFieldFromNodeWithDefaults(tmpl.Tree.Root, 0)
	if err != nil {
		return nil, fmt.Errorf("解析模板文件 %s 存在异常: %v", fileName, err)
	}

	return result, nil
}

// getFieldFromNode extracts variables from template nodes
func (p *Parser) getFieldFromNode(node parse.Node, depth int) ([]string, error) {
	depth = depth + 1
	if depth > maxDepth {
		return nil, errors.New("模板变量层级太深，检查模板是否正确")
	}
	var result []string
	switch node := node.(type) {
	case *parse.FieldNode:
		ident := node.Ident
		join := strings.Join(ident, ".")
		result = append(result, join)
	case *parse.CommandNode:
		args := node.Args
		firstWord := args[0]
		if firstWord.Type() == parse.NodeIdentifier {
			sonResult, err := p.parseCustomFunc(args, depth)
			if err != nil {
				return nil, err
			}
			result = append(result, sonResult...)
		} else {
			for _, arg := range args {
				sonResult, err := p.getFieldFromNode(arg, depth)
				if err != nil {
					return nil, err
				}
				result = append(result, sonResult...)
			}
		}
	case *parse.ActionNode:
		pipe := node.Pipe
		sonResult, err := p.getFieldFromNode(pipe, depth)
		if err != nil {
			return nil, err
		}
		result = append(result, sonResult...)

	case *parse.PipeNode:
		cmds := node.Cmds
		for _, cmd := range cmds {
			sonResult, err := p.getFieldFromNode(cmd, depth)
			if err != nil {
				return nil, err
			}
			result = append(result, sonResult...)
		}
	case *parse.ListNode:
		nodes := node.Nodes
		for _, item := range nodes {
			sonResult, err := p.getFieldFromNode(item, depth)
			if err != nil {
				return nil, err
			}
			result = append(result, sonResult...)
		}
	case *parse.IfNode:
		sonResult, err := p.processIfAndWithAndRange(node.Pipe, node.List, node.ElseList, depth)
		if err != nil {
			return nil, err
		}
		result = append(result, sonResult...)
	case *parse.RangeNode:
		sonResult, err := p.processIfAndWithAndRange(node.Pipe, node.List, node.ElseList, depth)
		if err != nil {
			return nil, err
		}
		result = append(result, sonResult...)

	case *parse.WithNode:
		sonResult, err := p.processIfAndWithAndRange(node.Pipe, node.List, node.ElseList, depth)
		if err != nil {
			return nil, err
		}
		result = append(result, sonResult...)
	case *parse.IdentifierNode:
	case *parse.TextNode:
	case *parse.DotNode:
	case *parse.StringNode:
	case *parse.BoolNode:
	case *parse.NilNode:
	case *parse.NumberNode:
	case *parse.VariableNode:
	}
	return result, nil
}

// getFieldFromNodeWithDefaults extracts variables with default values from template nodes
func (p *Parser) getFieldFromNodeWithDefaults(node parse.Node, depth int) ([]VariableInfo, error) {
	depth = depth + 1
	if depth > maxDepth {
		return nil, errors.New("模板变量层级太深，检查模板是否正确")
	}
	var result []VariableInfo
	switch node := node.(type) {
	case *parse.FieldNode:
		ident := node.Ident
		join := strings.Join(ident, ".")
		result = append(result, VariableInfo{Name: join})
	case *parse.CommandNode:
		args := node.Args
		firstWord := args[0]
		if firstWord.Type() == parse.NodeIdentifier {
			sonResult, err := p.parseCustomFuncWithDefaults(args, depth)
			if err != nil {
				return nil, err
			}
			result = append(result, sonResult...)
		} else {
			for _, arg := range args {
				sonResult, err := p.getFieldFromNodeWithDefaults(arg, depth)
				if err != nil {
					return nil, err
				}
				result = append(result, sonResult...)
			}
		}
	case *parse.ActionNode:
		pipe := node.Pipe
		sonResult, err := p.getFieldFromNodeWithDefaults(pipe, depth)
		if err != nil {
			return nil, err
		}
		result = append(result, sonResult...)

	case *parse.PipeNode:
		cmds := node.Cmds
		for _, cmd := range cmds {
			sonResult, err := p.getFieldFromNodeWithDefaults(cmd, depth)
			if err != nil {
				return nil, err
			}
			result = append(result, sonResult...)
		}
	case *parse.ListNode:
		nodes := node.Nodes
		for _, item := range nodes {
			sonResult, err := p.getFieldFromNodeWithDefaults(item, depth)
			if err != nil {
				return nil, err
			}
			result = append(result, sonResult...)
		}
	case *parse.IfNode:
		sonResult, err := p.processIfAndWithAndRangeWithDefaults(node.Pipe, node.List, node.ElseList, depth)
		if err != nil {
			return nil, err
		}
		result = append(result, sonResult...)
	case *parse.RangeNode:
		sonResult, err := p.processIfAndWithAndRangeWithDefaults(node.Pipe, node.List, node.ElseList, depth)
		if err != nil {
			return nil, err
		}
		result = append(result, sonResult...)

	case *parse.WithNode:
		sonResult, err := p.processIfAndWithAndRangeWithDefaults(node.Pipe, node.List, node.ElseList, depth)
		if err != nil {
			return nil, err
		}
		result = append(result, sonResult...)
	case *parse.IdentifierNode:
	case *parse.TextNode:
	case *parse.DotNode:
	case *parse.StringNode:
	case *parse.BoolNode:
	case *parse.NilNode:
	case *parse.NumberNode:
	case *parse.VariableNode:
	}
	return result, nil
}

// processIfAndWithAndRange processes if, range, and with nodes
func (p *Parser) processIfAndWithAndRange(pipe *parse.PipeNode, list, elseList *parse.ListNode, cycle int) ([]string, error) {
	var result []string
	sonResult, err := p.getFieldFromNode(pipe, cycle)
	if err != nil {
		return nil, err
	}
	result = append(result, sonResult...)
	sonResult, err = p.getFieldFromNode(list, cycle)
	if err != nil {
		return nil, err
	}
	result = append(result, sonResult...)
	if elseList != nil {
		sonResult, err = p.getFieldFromNode(elseList, cycle)
		if err != nil {
			return nil, err
		}
		result = append(result, sonResult...)
	}
	return result, nil
}

// processIfAndWithAndRangeWithDefaults processes if, range, and with nodes with default values
func (p *Parser) processIfAndWithAndRangeWithDefaults(pipe *parse.PipeNode, list, elseList *parse.ListNode, cycle int) ([]VariableInfo, error) {
	var result []VariableInfo
	sonResult, err := p.getFieldFromNodeWithDefaults(pipe, cycle)
	if err != nil {
		return nil, err
	}
	result = append(result, sonResult...)
	sonResult, err = p.getFieldFromNodeWithDefaults(list, cycle)
	if err != nil {
		return nil, err
	}
	result = append(result, sonResult...)
	if elseList != nil {
		sonResult, err = p.getFieldFromNodeWithDefaults(elseList, cycle)
		if err != nil {
			return nil, err
		}
		result = append(result, sonResult...)
	}
	return result, nil
}

// parseCustomFunc parses custom functions for variable extraction
func (p *Parser) parseCustomFunc(args []parse.Node, cycle int) ([]string, error) {
	var result []string
	node := args[0].(*parse.IdentifierNode)
	funcName := node.Ident
	match := p.functionMatcher.MatchCustomFunc(funcName)
	if match {
		if len(args) > 1 {
			item := args[1]
			if item.Type() == parse.NodeString {
				stringNode := item.(*parse.StringNode)
				result = append(result, stringNode.Text)
			} else {
				sonResult, err := p.getFieldFromNode(item, cycle)
				if err != nil {
					return nil, err
				}
				result = append(result, sonResult...)
			}
		}
	} else {
		for _, arg := range args {
			sonResult, err := p.getFieldFromNode(arg, cycle)
			if err != nil {
				return nil, err
			}
			result = append(result, sonResult...)
		}
	}
	return result, nil
}

// parseCustomFuncWithDefaults parses custom functions with default values
func (p *Parser) parseCustomFuncWithDefaults(args []parse.Node, cycle int) ([]VariableInfo, error) {
	var result []VariableInfo
	node := args[0].(*parse.IdentifierNode)
	funcName := node.Ident
	match := p.functionMatcher.MatchCustomFunc(funcName)
	if match {
		if len(args) > 1 {
			item := args[1]
			if item.Type() == parse.NodeString {
				stringNode := item.(*parse.StringNode)
				varInfo := VariableInfo{
					Name: stringNode.Text,
				}
				if len(args) > 2 && args[2].Type() == parse.NodeString {
					defaultValueNode := args[2].(*parse.StringNode)
					varInfo.DefaultValue = defaultValueNode.Text
				}
				result = append(result, varInfo)
			} else {
				sonResult, err := p.getFieldFromNode(item, cycle)
				if err != nil {
					return nil, err
				}
				for _, name := range sonResult {
					result = append(result, VariableInfo{Name: name})
				}
			}
		}
	} else {
		for _, arg := range args {
			sonResult, err := p.getFieldFromNode(arg, cycle)
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

// createMinimalFuncMap creates minimal function map for parsing
func (p *Parser) createMinimalFuncMap() template.FuncMap {
	return template.FuncMap{
		"getv":  func(key string, v ...string) string { return "" },
		"exists": func(key string) bool { return false },
		"get":   func(key string) (interface{}, error) { return nil, nil },
		"jsonv": func(key string) (map[string]interface{}, error) { return nil, nil },
	}
}