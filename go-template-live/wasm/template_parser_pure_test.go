//go:build !js
// +build !js

package main

import (
	"reflect"
	"testing"
)

// Test functions that don't require WASM/JS dependencies

func TestDefaultFunctionMatcher_MatchCustomFunc_Pure(t *testing.T) {
	matcher := &DefaultFunctionMatcher{}

	tests := []struct {
		name     string
		funcName string
		expected bool
	}{
		{
			name:     "matching function - exists",
			funcName: "exists",
			expected: true,
		},
		{
			name:     "matching function - get",
			funcName: "get",
			expected: true,
		},
		{
			name:     "matching function - getv",
			funcName: "getv",
			expected: true,
		},
		{
			name:     "matching function - jsonv",
			funcName: "jsonv",
			expected: true,
		},
		{
			name:     "non-matching function - printf",
			funcName: "printf",
			expected: false,
		},
		{
			name:     "non-matching function - eq",
			funcName: "eq",
			expected: false,
		},
		{
			name:     "empty function name",
			funcName: "",
			expected: false,
		},
		{
			name:     "case sensitive test",
			funcName: "Exists",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matcher.MatchCustomFunc(tt.funcName)
			if result != tt.expected {
				t.Errorf("MatchCustomFunc() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewFunctionHandler_Pure(t *testing.T) {
	variables := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}

	handler := NewFunctionHandler(variables)

	if handler == nil {
		t.Fatal("NewFunctionHandler() returned nil")
	}

	if !reflect.DeepEqual(handler.variables, variables) {
		t.Errorf("NewFunctionHandler() variables = %v, want %v", handler.variables, variables)
	}
}

func TestFunctionHandler_getvFunc_Pure(t *testing.T) {
	tests := []struct {
		name      string
		variables map[string]interface{}
		key       string
		defaults  []string
		expected  string
	}{
		{
			name: "existing key with non-empty value",
			variables: map[string]interface{}{
				"name": "john",
			},
			key:      "name",
			defaults: []string{"default"},
			expected: "john",
		},
		{
			name: "existing key with empty value",
			variables: map[string]interface{}{
				"name": "",
			},
			key:      "name",
			defaults: []string{"default"},
			expected: "default",
		},
		{
			name: "non-existing key with default",
			variables: map[string]interface{}{
				"other": "value",
			},
			key:      "missing",
			defaults: []string{"default"},
			expected: "default",
		},
		{
			name: "non-existing key without default",
			variables: map[string]interface{}{
				"other": "value",
			},
			key:      "missing",
			defaults: []string{},
			expected: "",
		},
		{
			name: "non-string value",
			variables: map[string]interface{}{
				"age": 25,
			},
			key:      "age",
			defaults: []string{"unknown"},
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewFunctionHandler(tt.variables)
			result := handler.getvFunc(tt.key, tt.defaults...)
			if result != tt.expected {
				t.Errorf("getvFunc() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFunctionHandler_existsFunc_Pure(t *testing.T) {
	variables := map[string]interface{}{
		"name": "john",
		"age":  25,
	}

	handler := NewFunctionHandler(variables)

	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{
			name:     "existing key",
			key:      "name",
			expected: true,
		},
		{
			name:     "existing key with numeric value",
			key:      "age",
			expected: true,
		},
		{
			name:     "non-existing key",
			key:      "missing",
			expected: false,
		},
		{
			name:     "empty key",
			key:      "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.existsFunc(tt.key)
			if result != tt.expected {
				t.Errorf("existsFunc() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFunctionHandler_getFunc_Pure(t *testing.T) {
	variables := map[string]interface{}{
		"name": "john",
		"age":  25,
	}

	handler := NewFunctionHandler(variables)

	tests := []struct {
		name      string
		key       string
		wantValue interface{}
		wantErr   bool
	}{
		{
			name:      "existing key",
			key:       "name",
			wantValue: "john",
			wantErr:   false,
		},
		{
			name:      "existing key with numeric value",
			key:       "age",
			wantValue: 25,
			wantErr:   false,
		},
		{
			name:      "non-existing key",
			key:       "missing",
			wantValue: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handler.getFunc(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantValue) {
				t.Errorf("getFunc() = %v, want %v", got, tt.wantValue)
			}
		})
	}
}

func TestNewParser_Pure(t *testing.T) {
	matcher := &DefaultFunctionMatcher{}

	parser := NewParser(matcher)

	if parser == nil {
		t.Fatal("NewParser() returned nil")
	}

	if parser.functionMatcher != matcher {
		t.Errorf("NewParser() functionMatcher = %v, want %v", parser.functionMatcher, matcher)
	}
}

func TestFunctionHandler_CreateFuncMap_Pure(t *testing.T) {
	handler := NewFunctionHandler(map[string]interface{}{
		"key1": "value1",
	})

	funcMap := handler.CreateFuncMap()

	if funcMap == nil {
		t.Fatal("CreateFuncMap() returned nil")
	}

	expectedFuncs := []string{"getv", "exists", "get", "jsonv"}
	for _, funcName := range expectedFuncs {
		if _, exists := funcMap[funcName]; !exists {
			t.Errorf("CreateFuncMap() missing function %s", funcName)
		}
	}
}