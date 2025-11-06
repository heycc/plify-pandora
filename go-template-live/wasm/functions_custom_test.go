//go:build !js && custom
// +build !js,custom

package main

import (
	"reflect"
	"strings"
	"testing"
	"text/template"
)

// createCustomParser creates a parser with custom functions enabled for testing
// Uses the actual production function registration from functions_custom_core.go
func createCustomParser() *Parser {
	registry := NewFunctionRegistry()

	// Save the current global registry
	oldRegistry := globalRegistry
	// Temporarily set our test registry as global
	globalRegistry = registry

	// Call registerCustomFunctions which is now available from functions_custom_core.go
	registerCustomFunctions()

	// Restore the old registry
	globalRegistry = oldRegistry

	return NewParser(registry)
}

// TestEndToEnd_CustomFunctions tests the complete workflow with custom functions:
// 1. Extract variables from template
// 2. Provide values for those variables
// 3. Render the template with those values
func TestEndToEnd_CustomFunctions(t *testing.T) {
	parserCustom := createCustomParser()

	tests := []struct {
		name           string
		template       string
		expectedVars   []VariableInfo
		providedValues map[string]interface{}
		expectedOutput string
	}{
		{
			name:     "simple field extraction and render",
			template: `Hello {{.Name}}!`,
			expectedVars: []VariableInfo{
				{Name: "Name"},
			},
			providedValues: map[string]interface{}{
				"Name": "John",
			},
			expectedOutput: "Hello John!",
		},
		{
			name:     "getv with default value - using provided value",
			template: `Username: {{getv "username" "guest"}}`,
			expectedVars: []VariableInfo{
				{Name: "username", DefaultValue: "guest"},
			},
			providedValues: map[string]interface{}{
				"username": "john_doe",
			},
			expectedOutput: "Username: john_doe",
		},
		{
			name:     "getv with default value - using default",
			template: `Username: {{getv "username" "guest"}}`,
			expectedVars: []VariableInfo{
				{Name: "username", DefaultValue: "guest"},
			},
			providedValues: map[string]interface{}{},
			expectedOutput: "Username: guest",
		},
		{
			name:     "exists function - key exists",
			template: `{{if exists "feature_flag"}}Feature enabled{{else}}Feature disabled{{end}}`,
			expectedVars: []VariableInfo{
				{Name: "feature_flag"},
			},
			providedValues: map[string]interface{}{
				"feature_flag": "true",
			},
			expectedOutput: "Feature enabled",
		},
		{
			name:     "exists function - key missing",
			template: `{{if exists "feature_flag"}}Feature enabled{{else}}Feature disabled{{end}}`,
			expectedVars: []VariableInfo{
				{Name: "feature_flag"},
			},
			providedValues: map[string]interface{}{},
			expectedOutput: "Feature disabled",
		},
		{
			name:     "mixed standard and custom functions",
			template: `Hello {{.Name}}, your username is {{getv "username" "anonymous"}} and email is {{.Email}}`,
			expectedVars: []VariableInfo{
				{Name: "Name"},
				{Name: "username", DefaultValue: "anonymous"},
				{Name: "Email"},
			},
			providedValues: map[string]interface{}{
				"Name":  "Alice",
				"Email": "alice@example.com",
			},
			expectedOutput: "Hello Alice, your username is anonymous and email is alice@example.com",
		},
		{
			name:     "multiple custom functions",
			template: `User: {{getv "name" "Unknown"}}, Active: {{exists "active"}}, Role: {{getv "role" "user"}}`,
			expectedVars: []VariableInfo{
				{Name: "name", DefaultValue: "Unknown"},
				{Name: "active"},
				{Name: "role", DefaultValue: "user"},
			},
			providedValues: map[string]interface{}{
				"name":   "Bob",
				"active": "yes",
			},
			expectedOutput: "User: Bob, Active: true, Role: user",
		},
		{
			name:     "json function - parse and access fields",
			template: `{{$user := json "user_data"}}Name: {{$user.name}}, Age: {{$user.age}}`,
			expectedVars: []VariableInfo{
				{Name: "user_data"},
			},
			providedValues: map[string]interface{}{
				"user_data": `{"name":"Alice","age":30}`,
			},
			expectedOutput: "Name: Alice, Age: 30",
		},
		{
			name:     "jsonArray function - iterate over array",
			template: `{{range jsonArray "fruits"}}{{.}}, {{end}}`,
			expectedVars: []VariableInfo{
				{Name: "fruits"},
			},
			providedValues: map[string]interface{}{
				"fruits": `["apple","banana","cherry"]`,
			},
			expectedOutput: "apple, banana, cherry, ",
		},
		{
			name:     "jsonArray function - iterate over number array",
			template: `{{range jsonArray "numbers"}}{{.}} {{end}}`,
			expectedVars: []VariableInfo{
				{Name: "numbers"},
			},
			providedValues: map[string]interface{}{
				"numbers": `[1,2,3,4,5]`,
			},
			expectedOutput: "1 2 3 4 5 ",
		},
		{
			name:     "jsonArray function - count items",
			template: `Count: {{len (jsonArray "items")}}`,
			expectedVars: []VariableInfo{
				{Name: "items"},
			},
			providedValues: map[string]interface{}{
				"items": `["a","b","c"]`,
			},
			expectedOutput: "Count: 3",
		},
		{
			name:     "jsonArray function - index access",
			template: `{{$arr := jsonArray "colors"}}First: {{index $arr 0}}, Second: {{index $arr 1}}`,
			expectedVars: []VariableInfo{
				{Name: "colors"},
			},
			providedValues: map[string]interface{}{
				"colors": `["red","green","blue"]`,
			},
			expectedOutput: "First: red, Second: green",
		},
		{
			name:     "combined json and jsonArray",
			template: `{{$config := json "config"}}App: {{$config.name}}, {{range jsonArray "tags"}}#{{.}} {{end}}`,
			expectedVars: []VariableInfo{
				{Name: "config"},
				{Name: "tags"},
			},
			providedValues: map[string]interface{}{
				"config": `{"name":"MyApp","version":"1.0"}`,
				"tags":   `["production","stable","v1"]`,
			},
			expectedOutput: "App: MyApp, #production #stable #v1 ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Step 1: Extract variables
			extractedVars, err := parserCustom.ExtractVariablesWithDefaults("test.tmpl", tt.template)
			if err != nil {
				t.Fatalf("ExtractVariablesWithDefaults() error = %v", err)
			}
			// Verify extracted variables match expected
			if !reflect.DeepEqual(extractedVars, tt.expectedVars) {
				t.Errorf("ExtractVariablesWithDefaults() = %v, want %v", extractedVars, tt.expectedVars)
			}

			// Step 2: Render template with provided values
			rendered, err := renderTemplateWithCustomFunctions(tt.template, tt.providedValues)
			if err != nil {
				t.Fatalf("renderTemplateWithCustomFunctions() error = %v", err)
			}
			// Verify rendered output
			if rendered != tt.expectedOutput {
				t.Errorf("renderTemplateWithCustomFunctions() = %q, want %q", rendered, tt.expectedOutput)
			}
		})
	}
}

// Helper function for rendering templates in tests

// renderTemplateWithCustomFunctions renders a template with custom functions enabled
// Uses the actual production implementation from functions_custom_core.go
func renderTemplateWithCustomFunctions(templateContent string, variables map[string]interface{}) (string, error) {
	// Use the actual production implementation - this is what we're testing!
	funcMap := GetCustomRenderFuncMap(variables)

	tmpl, err := template.New("test").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	err = tmpl.Execute(&result, variables)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}
