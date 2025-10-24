package main

import (
	"text/template"
	"text/template/parse"
)

const maxDepth = 40

// VariableInfo stores variable information including name and default value
type VariableInfo struct {
	Name         string `json:"name"`
	DefaultValue string `json:"defaultValue,omitempty"`
}

// VariableExtractor extracts variable names from function arguments
// This allows us to identify which variables are being used in templates
// similar to how Confd tracks template dependencies
type VariableExtractor func(args []parse.Node, cycle int) ([]string, error)

// VariableExtractorWithDefaults extracts variables with default values
type VariableExtractorWithDefaults func(args []parse.Node, cycle int) ([]VariableInfo, error)

// FunctionDefinition defines a custom template function
// Inspired by Confd's approach to template functions
// https://github.com/kelseyhightower/confd
type FunctionDefinition struct {
	Name                  string
	Description           string
	Handler               interface{}
	Extractor             VariableExtractor
	ExtractorWithDefaults VariableExtractorWithDefaults
}

// FunctionRegistry manages all custom template functions
// This is the single source of truth for available functions
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

// HasFunction checks if a function is registered
func (r *FunctionRegistry) HasFunction(name string) bool {
	_, exists := r.functions[name]
	return exists
}

// GetFunctionNames returns all registered function names
func (r *FunctionRegistry) GetFunctionNames() []string {
	names := make([]string, 0, len(r.functions))
	for name := range r.functions {
		names = append(names, name)
	}
	return names
}

// GetMinimalFuncMap creates a minimal function map for parsing
// Uses minimal implementations that don't require actual variables
func (r *FunctionRegistry) GetMinimalFuncMap() template.FuncMap {
	funcMap := template.FuncMap{}
	for name, def := range r.functions {
		funcMap[name] = def.Handler
	}
	return funcMap
}

// GetRenderFuncMap creates a function map for rendering with actual variable values
func (r *FunctionRegistry) GetRenderFuncMap(variables map[string]interface{}) template.FuncMap {
	funcMap := template.FuncMap{}
	for name, def := range r.functions {
		// For rendering, we need to create closures with the variables
		// This is handled by the function implementations themselves
		funcMap[name] = def.Handler
	}
	return funcMap
}
