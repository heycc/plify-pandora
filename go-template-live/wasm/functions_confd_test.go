//go:build !js && confd
// +build !js,confd

package main

import (
	"reflect"
	"strings"
	"testing"
	"text/template"
)

// createConfdParser creates a parser with Confd-style functions enabled for testing
// Uses the actual production function registration from functions_confd.go
func createConfdParser() *Parser {
	// Register Confd functions in the global registry
	// Note: functions_confd.go has init() that calls this, but that file has 'js && confd' tag
	// so it won't be included in tests (!js && confd). We call it manually here.
	registerConfdFunctions()
	return NewParser(GetGlobalRegistry())
}

// TestEndToEnd_ConfdFunctions tests the complete workflow with Confd-style functions:
// 1. Extract variables from template
// 2. Provide values for those variables
// 3. Render the template with those values
func TestEndToEnd_ConfdFunctions(t *testing.T) {
	parserConfd := createConfdParser()

	tests := []struct {
		name           string
		template       string
		expectedVars   []VariableInfo
		providedValues map[string]interface{}
		expectedOutput string
	}{
		// String manipulation functions
		{
			name:           "base function",
			template:       `Path: {{base "/home/user/file.txt"}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "Path: file.txt",
		},
		{
			name:           "base function with variable",
			template:       `Path: {{base .path}}`,
			expectedVars:   []VariableInfo{{Name: "path"}},
			providedValues: map[string]interface{}{"path": "/foo/bar/abc.txt"},
			expectedOutput: "Path: abc.txt",
		},
		{
			name:           "split function",
			template:       `{{range split "a,b,c" ","}}{{.}} {{end}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "a b c ",
		},
		{
			name:           "split function with variable",
			template:       `{{range split .csv ","}}{{.}} {{end}}`,
			expectedVars:   []VariableInfo{{Name: "csv"}},
			providedValues: map[string]interface{}{"csv": "x,y,z"},
			expectedOutput: "x y z ",
		},
		{
			name:           "dir function",
			template:       `Directory: {{dir "/home/user/file.txt"}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "Directory: /home/user",
		},
		{
			name:           "join function",
			template:       `{{join (split "a,b,c" ",") "-"}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "a-b-c",
		},
		{
			name:           "toUpper function",
			template:       `{{toUpper "hello world"}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "HELLO WORLD",
		},
		{
			name:           "toLower function",
			template:       `{{toLower "HELLO WORLD"}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "hello world",
		},
		{
			name:           "replace function",
			template:       `{{replace "hello world" "world" "golang" -1}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "hello golang",
		},
		{
			name:           "contains function - true",
			template:       `{{if contains "hello world" "world"}}Found{{else}}Not found{{end}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "Found",
		},
		{
			name:           "contains function - false",
			template:       `{{if contains "hello world" "golang"}}Found{{else}}Not found{{end}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "Not found",
		},
		{
			name:           "trimSuffix function",
			template:       `{{trimSuffix "filename.txt" ".txt"}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "filename",
		},

		// Encoding functions
		{
			name:           "base64Encode function",
			template:       `{{base64Encode "hello"}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "aGVsbG8=",
		},
		{
			name:           "base64Decode function",
			template:       `{{base64Decode "aGVsbG8="}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "hello",
		},

		// Boolean functions
		{
			name:           "parseBool function - true",
			template:       `{{if parseBool "true"}}Yes{{else}}No{{end}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "Yes",
		},
		{
			name:           "parseBool function - false",
			template:       `{{if parseBool "false"}}Yes{{else}}No{{end}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "No",
		},

		// Math functions
		{
			name:           "add function",
			template:       `{{add 5 3}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "8",
		},
		{
			name:           "sub function",
			template:       `{{sub 10 4}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "6",
		},
		{
			name:           "mul function",
			template:       `{{mul 6 7}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "42",
		},
		{
			name:           "div function",
			template:       `{{div 15 3}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "5",
		},
		{
			name:           "mod function",
			template:       `{{mod 17 5}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "2",
		},

		// Sequence functions
		{
			name:           "seq function",
			template:       `{{range seq 1 3}}{{.}} {{end}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "1 2 3 ",
		},

		// Conversion functions
		{
			name:           "atoi function",
			template:       `{{atoi "42"}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "42",
		},

		// Map functions
		{
			name:           "map function",
			template:       `{{$m := map "name" "Alice" "age" 30}}Name: {{$m.name}}, Age: {{$m.age}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "Name: Alice, Age: 30",
		},

		// Reverse function
		{
			name:           "reverse function with strings",
			template:       `{{range reverse (split "a,b,c" ",")}}{{.}} {{end}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "c b a ",
		},

		// JSON functions (these DO extract variables)
		{
			name:     "json function with variable",
			template: `{{$config := json "config_data"}}App: {{$config.name}}`,
			expectedVars: []VariableInfo{
				{Name: "config_data"},
			},
			providedValues: map[string]interface{}{
				"config_data": `{"name":"MyApp","version":"1.0"}`,
			},
			expectedOutput: "App: MyApp",
		},
		{
			name:     "jsonArray function with variable",
			template: `{{range jsonArray "items"}}{{.}} {{end}}`,
			expectedVars: []VariableInfo{
				{Name: "items"},
			},
			providedValues: map[string]interface{}{
				"items": `["apple","banana","cherry"]`,
			},
			expectedOutput: "apple banana cherry ",
		},

		// Combined functions
		{
			name:           "combined string functions",
			template:       `{{toUpper (base "/path/to/file.txt")}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "FILE.TXT",
		},
		{
			name:           "combined math and string functions",
			template:       `Result: {{add (atoi "10") (atoi "20")}}`,
			expectedVars:   nil,
			providedValues: map[string]interface{}{},
			expectedOutput: "Result: 30",
		},
		{
			name:     "mixed standard fields and Confd functions",
			template: `Hello {{.Name}}, your file is {{base .FilePath}} and total is {{add .Count 5}}`,
			expectedVars: []VariableInfo{
				{Name: "Name"},     // Standard Go template field
				{Name: "FilePath"}, // Extracted by base function
				{Name: "Count"},    // Extracted by add function
			},
			providedValues: map[string]interface{}{
				"Name":     "Bob",
				"FilePath": "/home/bob/document.pdf",
				"Count":    10,
			},
			expectedOutput: "Hello Bob, your file is document.pdf and total is 15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Step 1: Extract variables
			extractedVars, err := parserConfd.ExtractVariablesWithDefaults("test.tmpl", tt.template)
			if err != nil {
				t.Fatalf("ExtractVariablesWithDefaults() error = %v", err)
			}
			// Verify extracted variables match expected
			if !reflect.DeepEqual(extractedVars, tt.expectedVars) {
				t.Errorf("ExtractVariablesWithDefaults() = %v, want %v", extractedVars, tt.expectedVars)
			}

			// Step 2: Render template with provided values
			rendered, err := renderTemplateWithConfdFunctions(tt.template, tt.providedValues)
			if err != nil {
				t.Fatalf("renderTemplateWithConfdFunctions() error = %v", err)
			}
			// Verify rendered output
			if rendered != tt.expectedOutput {
				t.Errorf("renderTemplateWithConfdFunctions() = %q, want %q", rendered, tt.expectedOutput)
			}
		})
	}
}

// TestConfdFunctions_VariableExtraction tests that Confd functions correctly extract variables
func TestConfdFunctions_VariableExtraction(t *testing.T) {
	parserConfd := createConfdParser()

	tests := []struct {
		name         string
		template     string
		expectedVars []VariableInfo
	}{
		{
			name:         "json function extracts variable",
			template:     `{{json "my_data"}}`,
			expectedVars: []VariableInfo{{Name: "my_data"}},
		},
		{
			name:         "jsonArray function extracts variable",
			template:     `{{jsonArray "items"}}`,
			expectedVars: []VariableInfo{{Name: "items"}},
		},
		{
			name:         "pure utility functions with literals extract no variables",
			template:     `{{join (split "a,b" ",") "-"}} {{map "key" "value"}}`,
			expectedVars: nil,
		},
		{
			name:         "string transformation functions with literals extract no variables",
			template:     `{{base "path"}} {{split "a,b" ","}} {{dir "path"}}`,
			expectedVars: nil,
		},
		{
			name:         "case transformation functions with literals extract no variables",
			template:     `{{toUpper "hello"}} {{toLower "WORLD"}}`,
			expectedVars: nil,
		},
		{
			name:         "string functions with literals extract no variables",
			template:     `{{contains "test" "es"}} {{replace "hello" "l" "L" -1}} {{trimSuffix "file.txt" ".txt"}}`,
			expectedVars: nil,
		},
		{
			name:         "encoding functions with literals extract no variables",
			template:     `{{base64Encode "hello"}} {{base64Decode "aGVsbG8="}}`,
			expectedVars: nil,
		},
		{
			name:         "math functions with literals extract no variables",
			template:     `{{add 1 2}} {{sub 5 3}} {{mul 2 3}} {{div 10 2}} {{mod 7 3}}`,
			expectedVars: nil,
		},
		{
			name:         "seq function with literals extracts no variables",
			template:     `{{range seq 1 5}}{{.}}{{end}}`,
			expectedVars: nil,
		},
		{
			name:         "string transformation functions with field access extract them",
			template:     `{{base .filepath}} {{dir .filepath}} {{toUpper .name}}`,
			expectedVars: []VariableInfo{{Name: "filepath"}, {Name: "filepath"}, {Name: "name"}},
		},
		{
			name:         "math functions with field access extract them",
			template:     `{{add .count 5}} {{sub .total 10}} {{mul .value 2}}`,
			expectedVars: []VariableInfo{{Name: "count"}, {Name: "total"}, {Name: "value"}},
		},
		{
			name:         "string functions with field access extract them",
			template:     `{{toLower .title}} {{contains .text "search"}} {{replace .content "old" "new" -1}}`,
			expectedVars: []VariableInfo{{Name: "title"}, {Name: "text"}, {Name: "content"}},
		},
		{
			name:         "parseBool function with field access extracts variable",
			template:     `{{parseBool .enabled}}`,
			expectedVars: []VariableInfo{{Name: "enabled"}},
		},
		{
			name:         "base64 functions with field access extract variables",
			template:     `{{base64Encode .data}} {{base64Decode .encoded}}`,
			expectedVars: []VariableInfo{{Name: "data"}, {Name: "encoded"}},
		},
		{
			name:         "mixed json and utility with literals",
			template:     `{{json "data"}} {{base "path"}}`,
			expectedVars: []VariableInfo{{Name: "data"}}, // base with literal doesn't extract
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractedVars, err := parserConfd.ExtractVariablesWithDefaults("test.tmpl", tt.template)
			if err != nil {
				t.Fatalf("ExtractVariablesWithDefaults() error = %v", err)
			}
			if !reflect.DeepEqual(extractedVars, tt.expectedVars) {
				t.Errorf("ExtractVariablesWithDefaults() = %v, want %v", extractedVars, tt.expectedVars)
			}
		})
	}
}

// Helper functions for rendering templates in tests

// renderTemplateWithConfdFunctions renders a template with Confd functions enabled
// Uses the actual production implementation from functions_confd_core.go
func renderTemplateWithConfdFunctions(templateContent string, variables map[string]interface{}) (string, error) {
	// Use the actual production implementation - this is what we're testing!
	funcMap := GetConfdRenderFuncMap(variables)

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
