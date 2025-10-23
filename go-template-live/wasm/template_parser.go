package main

import (
	"text/template"
)

const maxDepth = 40

// VariableInfo stores variable information including name and default value
type VariableInfo struct {
	Name         string `json:"name"`
	DefaultValue string `json:"defaultValue"`
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
		"getv":   getvFunc(h.variables),
		"exists": existsFunc(h.variables),
		"get":    getFunc(h.variables),
		"jsonv":  jsonvFunc(h.variables),
	}
}
