//go:build !js && official
// +build !js,official

package main

import (
	"reflect"
	"strings"
	"testing"
	"text/template"
)

// TestHelper provides utilities for testing different build configurations
type TestHelper struct{}

// NewTestHelper creates a new test helper
func NewTestHelper() *TestHelper {
	return &TestHelper{}
}

// NewParserWithOfficialFunctions creates a parser with only official functions
// This simulates the official build with no custom functions registered
func (h *TestHelper) NewParserWithOfficialFunctions() *Parser {
	registry := NewFunctionRegistry()
	// Don't register any custom functions
	return NewParser(registry)
}

// TestEndToEnd_OfficialMode tests that official mode works correctly
// This tests standard Go template functionality without any custom functions
func TestEndToEnd_OfficialMode(t *testing.T) {
	helper := NewTestHelper()
	parserOfficial := helper.NewParserWithOfficialFunctions()

	tests := []struct {
		name           string
		template       string
		expectedVars   []VariableInfo
		providedValues map[string]interface{}
		expectedOutput string
	}{
		{
			name:     "standard template fields only",
			template: `Hello {{.Name}}, your email is {{.Email}}`,
			expectedVars: []VariableInfo{
				{Name: "Name"},
				{Name: "Email"},
			},
			providedValues: map[string]interface{}{
				"Name":  "Charlie",
				"Email": "charlie@example.com",
			},
			expectedOutput: "Hello Charlie, your email is charlie@example.com",
		},
		{
			name:     "with if statement",
			template: `{{if .Active}}User is active{{else}}User is inactive{{end}}`,
			expectedVars: []VariableInfo{
				{Name: "Active"},
			},
			providedValues: map[string]interface{}{
				"Active": true,
			},
			expectedOutput: "User is active",
		},
		{
			name:     "nested field access",
			template: `Hello {{.User.Name}}, your role is {{.User.Role}}`,
			expectedVars: []VariableInfo{
				{Name: "User.Name"},
				{Name: "User.Role"},
			},
			providedValues: map[string]interface{}{
				"User": map[string]interface{}{
					"Name": "David",
					"Role": "Admin",
				},
			},
			expectedOutput: "Hello David, your role is Admin",
		},
		{
			name:     "range over slice",
			template: `{{range .Items}}{{.}}, {{end}}`,
			expectedVars: []VariableInfo{
				{Name: "Items"},
			},
			providedValues: map[string]interface{}{
				"Items": []string{"apple", "banana", "cherry"},
			},
			expectedOutput: "apple, banana, cherry, ",
		},
		{
			name:     "with statement",
			template: `{{with .User}}Name: {{.Name}}{{end}}`,
			expectedVars: []VariableInfo{
				{Name: "User"},
				{Name: "Name"},
			},
			providedValues: map[string]interface{}{
				"User": map[string]interface{}{
					"Name": "Eve",
				},
			},
			expectedOutput: "Name: Eve",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Step 1: Extract variables
			extractedVars, err := parserOfficial.ExtractVariablesWithDefaults("test.tmpl", tt.template)
			if err != nil {
				t.Fatalf("ExtractVariablesWithDefaults() error = %v", err)
			}
			// Verify extracted variables match expected
			if !reflect.DeepEqual(extractedVars, tt.expectedVars) {
				t.Errorf("ExtractVariablesWithDefaults() = %v, want %v", extractedVars, tt.expectedVars)
			}

			// Step 2: Render template (without custom functions)
			rendered, err := renderTemplateOfficialMode(tt.template, tt.providedValues)
			if err != nil {
				t.Fatalf("renderTemplateOfficialMode() error = %v", err)
			}
			// Verify rendered output
			if rendered != tt.expectedOutput {
				t.Errorf("renderTemplateOfficialMode() = %q, want %q", rendered, tt.expectedOutput)
			}
		})
	}
}

// renderTemplateOfficialMode renders a template without custom functions (official mode)
func renderTemplateOfficialMode(templateContent string, variables map[string]interface{}) (string, error) {
	tmpl, err := template.New("test").Parse(templateContent)
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
