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
	return h.CreateFuncMapWithConfig(DefaultBuildConfig())
}

// CreateFuncMapWithConfig creates template function map based on build configuration
func (h *FunctionHandler) CreateFuncMapWithConfig(config *BuildConfig) template.FuncMap {
	funcMap := template.FuncMap{}

	if config.IncludeCustomFunctions {
		funcMap["getv"] = getvFunc(h.variables)
		funcMap["exists"] = existsFunc(h.variables)
		funcMap["get"] = getFunc(h.variables)
		funcMap["jsonv"] = jsonvFunc(h.variables)
	}

	return funcMap
}
