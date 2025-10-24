//go:build !js
// +build !js

package main

import (
	"reflect"
	"testing"
)

// Test functions that don't require WASM/JS dependencies

func TestFunctionRegistry_RegisterFunction_Pure(t *testing.T) {
	registry := NewFunctionRegistry()

	funcDef := &FunctionDefinition{
		Name:        "testfunc",
		Description: "Test function",
		Handler:     func() string { return "test" },
	}

	registry.RegisterFunction(funcDef)

	got, exists := registry.GetFunction("testfunc")
	if !exists {
		t.Error("RegisterFunction() function not found after registration")
	}
	if got != funcDef {
		t.Errorf("RegisterFunction() = %v, want %v", got, funcDef)
	}
}

func TestFunctionRegistry_HasFunction_Pure(t *testing.T) {
	registry := NewFunctionRegistry()

	if registry.HasFunction("testfunc") {
		t.Error("HasFunction() returned true for non-existent function")
	}

	registry.RegisterFunction(&FunctionDefinition{
		Name:    "testfunc",
		Handler: func() string { return "test" },
	})

	if !registry.HasFunction("testfunc") {
		t.Error("HasFunction() returned false for registered function")
	}
}

func TestFunctionRegistry_GetFunctionNames_Pure(t *testing.T) {
	registry := NewFunctionRegistry()

	registry.RegisterFunction(&FunctionDefinition{
		Name:    "func1",
		Handler: func() string { return "test1" },
	})
	registry.RegisterFunction(&FunctionDefinition{
		Name:    "func2",
		Handler: func() string { return "test2" },
	})

	names := registry.GetFunctionNames()
	if len(names) != 2 {
		t.Errorf("GetFunctionNames() returned %d functions, want 2", len(names))
	}

	// Check that both function names are present
	nameMap := make(map[string]bool)
	for _, name := range names {
		nameMap[name] = true
	}
	if !nameMap["func1"] || !nameMap["func2"] {
		t.Errorf("GetFunctionNames() = %v, want [func1, func2]", names)
	}
}

func TestGetvRenderHandler_Pure(t *testing.T) {
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
			result := getvRenderHandler(tt.variables)(tt.key, tt.defaults...)
			if result != tt.expected {
				t.Errorf("getvRenderHandler() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExistsRenderHandler_Pure(t *testing.T) {
	variables := map[string]interface{}{
		"name": "john",
		"age":  25,
	}

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
			result := existsRenderHandler(variables)(tt.key)
			if result != tt.expected {
				t.Errorf("existsRenderHandler() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetRenderHandler_Pure(t *testing.T) {
	variables := map[string]interface{}{
		"name": "john",
		"age":  25,
	}

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
			got, err := getRenderHandler(variables)(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRenderHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantValue) {
				t.Errorf("getRenderHandler() = %v, want %v", got, tt.wantValue)
			}
		})
	}
}

func TestJsonRenderHandler_Pure(t *testing.T) {
	variables := map[string]interface{}{
		"validJson":   `{"name":"john","age":25}`,
		"invalidJson": `{invalid}`,
		"notString":   123,
	}

	tests := []struct {
		name      string
		key       string
		wantValue map[string]interface{}
		wantErr   bool
	}{
		{
			name: "valid json",
			key:  "validJson",
			wantValue: map[string]interface{}{
				"name": "john",
				"age":  float64(25),
			},
			wantErr: false,
		},
		{
			name:      "invalid json",
			key:       "invalidJson",
			wantValue: nil,
			wantErr:   true,
		},
		{
			name:      "non-existing key",
			key:       "missing",
			wantValue: nil,
			wantErr:   true,
		},
		{
			name:      "non-string value",
			key:       "notString",
			wantValue: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonRenderHandler(variables)(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonRenderHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantValue) {
				t.Errorf("jsonRenderHandler() = %v, want %v", got, tt.wantValue)
			}
		})
	}
}

func TestJsonArrayRenderHandler_Pure(t *testing.T) {
	variables := map[string]interface{}{
		"validArray":   `["apple","banana","cherry"]`,
		"validNumbers": `[1,2,3,4,5]`,
		"validMixed":   `[1,"two",3.5,true,null]`,
		"invalidJson":  `[invalid]`,
		"notString":    123,
		"jsonObject":   `{"name":"john"}`,
	}

	tests := []struct {
		name      string
		key       string
		wantValue []interface{}
		wantErr   bool
	}{
		{
			name:      "valid string array",
			key:       "validArray",
			wantValue: []interface{}{"apple", "banana", "cherry"},
			wantErr:   false,
		},
		{
			name:      "valid number array",
			key:       "validNumbers",
			wantValue: []interface{}{float64(1), float64(2), float64(3), float64(4), float64(5)},
			wantErr:   false,
		},
		{
			name:      "valid mixed array",
			key:       "validMixed",
			wantValue: []interface{}{float64(1), "two", float64(3.5), true, nil},
			wantErr:   false,
		},
		{
			name:      "invalid json",
			key:       "invalidJson",
			wantValue: nil,
			wantErr:   true,
		},
		{
			name:      "non-existing key",
			key:       "missing",
			wantValue: nil,
			wantErr:   true,
		},
		{
			name:      "non-string value",
			key:       "notString",
			wantValue: nil,
			wantErr:   true,
		},
		{
			name:      "json object instead of array",
			key:       "jsonObject",
			wantValue: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonArrayRenderHandler(variables)(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonArrayRenderHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantValue) {
				t.Errorf("jsonArrayRenderHandler() = %v, want %v", got, tt.wantValue)
			}
		})
	}
}

func TestNewParser_Pure(t *testing.T) {
	registry := NewFunctionRegistry()

	parser := NewParser(registry)

	if parser == nil {
		t.Fatal("NewParser() returned nil")
	}

	if parser.registry != registry {
		t.Errorf("NewParser() registry = %v, want %v", parser.registry, registry)
	}
}

func TestParser_ExtractVariables_Pure(t *testing.T) {
	helper := NewTestHelper()

	// Test with custom functions
	parserCustom := helper.NewParserWithCustomFunctions()
	parserOfficial := helper.NewParserWithOfficialFunctions()

	tests := []struct {
		name             string
		fileName         string
		fileContent      string
		expectedCustom   []string
		expectedOfficial []string
		wantErrCustom    bool
		wantErrOfficial  bool
	}{
		{
			name:             "simple field",
			fileName:         "test.tmpl",
			fileContent:      `Hello {{.Name}}!`,
			expectedCustom:   []string{"Name"},
			expectedOfficial: []string{"Name"},
			wantErrCustom:    false,
			wantErrOfficial:  false,
		},
		{
			name:             "nested field",
			fileName:         "test.tmpl",
			fileContent:      `Hello {{.User.Name}}!`,
			expectedCustom:   []string{"User.Name"},
			expectedOfficial: []string{"User.Name"},
			wantErrCustom:    false,
			wantErrOfficial:  false,
		},
		{
			name:             "custom function with string parameter",
			fileName:         "test.tmpl",
			fileContent:      `{{getv "username"}}`,
			expectedCustom:   []string{"username"},
			expectedOfficial: nil,
			wantErrCustom:    false,
			wantErrOfficial:  true, // Official parser errors on undefined functions
		},
		{
			name:             "custom function with field parameter",
			fileName:         "test.tmpl",
			fileContent:      `{{getv .Key}}`,
			expectedCustom:   []string{"Key"},
			expectedOfficial: nil,
			wantErrCustom:    false,
			wantErrOfficial:  true, // Official parser errors on undefined functions
		},
		{
			name:             "multiple variables",
			fileName:         "test.tmpl",
			fileContent:      `Hello {{.Name}}, your email is {{.Email}} and username is {{getv "username"}}`,
			expectedCustom:   []string{"Name", "Email", "username"},
			expectedOfficial: nil,
			wantErrCustom:    false,
			wantErrOfficial:  true, // Official parser errors on undefined functions
		},
		{
			name:             "if statement",
			fileName:         "test.tmpl",
			fileContent:      `{{if .Enabled}}{{.Name}}{{end}}`,
			expectedCustom:   []string{"Enabled", "Name"},
			expectedOfficial: []string{"Enabled", "Name"},
			wantErrCustom:    false,
			wantErrOfficial:  false,
		},
		{
			name:             "range statement",
			fileName:         "test.tmpl",
			fileContent:      `{{range .Items}}{{.Name}}{{end}}`,
			expectedCustom:   []string{"Items", "Name"},
			expectedOfficial: []string{"Items", "Name"},
			wantErrCustom:    false,
			wantErrOfficial:  false,
		},
		{
			name:             "with statement",
			fileName:         "test.tmpl",
			fileContent:      `{{with .User}}{{.Name}}{{end}}`,
			expectedCustom:   []string{"User", "Name"},
			expectedOfficial: []string{"User", "Name"},
			wantErrCustom:    false,
			wantErrOfficial:  false,
		},
		{
			name:             "non-matching function should process arguments",
			fileName:         "test.tmpl",
			fileContent:      `{{printf "%s" "hello"}}`,
			expectedCustom:   []string{},
			expectedOfficial: []string{},
			wantErrCustom:    false,
			wantErrOfficial:  false,
		},
		{
			name:             "invalid template",
			fileName:         "test.tmpl",
			fileContent:      `{{.unclosed`,
			expectedCustom:   nil,
			expectedOfficial: nil,
			wantErrCustom:    true,
			wantErrOfficial:  true,
		},
		{
			name:             "template with else branch",
			fileName:         "test.tmpl",
			fileContent:      `{{if .Active}}{{.Name}}{{else}}{{.DefaultName}}{{end}}`,
			expectedCustom:   []string{"Active", "Name", "DefaultName"},
			expectedOfficial: []string{"Active", "Name", "DefaultName"},
			wantErrCustom:    false,
			wantErrOfficial:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_custom", func(t *testing.T) {
			got, err := parserCustom.ExtractVariables(tt.fileName, tt.fileContent)
			if (err != nil) != tt.wantErrCustom {
				t.Errorf("ExtractVariables() error = %v, wantErr %v", err, tt.wantErrCustom)
				return
			}
			if !tt.wantErrCustom {
				// Handle empty slice comparison
				if len(got) == 0 && len(tt.expectedCustom) == 0 {
					// Both empty, this is fine
				} else if !reflect.DeepEqual(got, tt.expectedCustom) {
					t.Errorf("ExtractVariables() = %v, want %v", got, tt.expectedCustom)
				}
			}
		})

		t.Run(tt.name+"_official", func(t *testing.T) {
			got, err := parserOfficial.ExtractVariables(tt.fileName, tt.fileContent)
			if (err != nil) != tt.wantErrOfficial {
				t.Errorf("ExtractVariables() error = %v, wantErr %v", err, tt.wantErrOfficial)
				return
			}
			if !tt.wantErrOfficial {
				// Handle empty slice comparison
				if len(got) == 0 && len(tt.expectedOfficial) == 0 {
					// Both empty, this is fine
				} else if !reflect.DeepEqual(got, tt.expectedOfficial) {
					t.Errorf("ExtractVariables() = %v, want %v", got, tt.expectedOfficial)
				}
			}
		})
	}
}
