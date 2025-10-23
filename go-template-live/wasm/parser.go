package main

import (
	"errors"
	"fmt"
	"strings"
	"text/template"
	"text/template/parse"
)

// Parser handles template parsing and variable extraction
// This follows the Confd pattern of parsing templates to extract dependencies
type Parser struct {
	functionMatcher  FunctionMatcher
	functionRegistry *FunctionRegistry
}

// NewParser creates a new template parser
func NewParser(matcher FunctionMatcher) *Parser {
	return &Parser{
		functionMatcher:  matcher,
		functionRegistry: DefaultFunctions(),
	}
}

// NewParserWithRegistry creates a new template parser with custom function registry
func NewParserWithRegistry(matcher FunctionMatcher, registry *FunctionRegistry) *Parser {
	return &Parser{
		functionMatcher:  matcher,
		functionRegistry: registry,
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
	return p.functionRegistry.GetMinimalFuncMap()
}
