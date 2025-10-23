//go:build !js
// +build !js

package main

import (
	"reflect"
	"testing"
)

func TestParser_ExtractVariables_Pure(t *testing.T) {
	matcher := &DefaultFunctionMatcher{}
	parser := NewParser(matcher)

	tests := []struct {
		name        string
		fileName    string
		fileContent string
		expected    []string
		wantErr     bool
	}{
		{
			name:        "simple field",
			fileName:    "test.tmpl",
			fileContent: `Hello {{.Name}}!`,
			expected:    []string{"Name"},
			wantErr:     false,
		},
		{
			name:        "nested field",
			fileName:    "test.tmpl",
			fileContent: `Hello {{.User.Name}}!`,
			expected:    []string{"User.Name"},
			wantErr:     false,
		},
		{
			name:        "custom function with string parameter",
			fileName:    "test.tmpl",
			fileContent: `{{getv "username"}}`,
			expected:    []string{"username"},
			wantErr:     false,
		},
		{
			name:        "custom function with field parameter",
			fileName:    "test.tmpl",
			fileContent: `{{getv .Key}}`,
			expected:    []string{"Key"},
			wantErr:     false,
		},
		{
			name:        "multiple variables",
			fileName:    "test.tmpl",
			fileContent: `Hello {{.Name}}, your email is {{.Email}} and {{getv "username"}}`,
			expected:    []string{"Name", "Email", "username"},
			wantErr:     false,
		},
		{
			name:        "if statement",
			fileName:    "test.tmpl",
			fileContent: `{{if .Enabled}}{{.Name}}{{end}}`,
			expected:    []string{"Enabled", "Name"},
			wantErr:     false,
		},
		{
			name:        "range statement",
			fileName:    "test.tmpl",
			fileContent: `{{range .Items}}{{.Name}}{{end}}`,
			expected:    []string{"Items", "Name"},
			wantErr:     false,
		},
		{
			name:        "with statement",
			fileName:    "test.tmpl",
			fileContent: `{{with .User}}{{.Name}}{{end}}`,
			expected:    []string{"User", "Name"},
			wantErr:     false,
		},
		{
			name:        "non-matching function should process arguments",
			fileName:    "test.tmpl",
			fileContent: `{{printf "%s" "hello"}}`,
			expected:    []string{},
			wantErr:     false,
		},
		{
			name:        "invalid template",
			fileName:    "test.tmpl",
			fileContent: `{{.unclosed`,
			expected:    nil,
			wantErr:     true,
		},
		{
			name:        "template with else branch",
			fileName:    "test.tmpl",
			fileContent: `{{if .Active}}{{.Name}}{{else}}{{.DefaultName}}{{end}}`,
			expected:    []string{"Active", "Name", "DefaultName"},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ExtractVariables(tt.fileName, tt.fileContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractVariables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Handle empty slice comparison
				if len(got) == 0 && len(tt.expected) == 0 {
					// Both empty, this is fine
				} else if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("ExtractVariables() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

func TestParser_ExtractVariablesWithDefaults_Pure(t *testing.T) {
	matcher := &DefaultFunctionMatcher{}
	parser := NewParser(matcher)

	tests := []struct {
		name        string
		fileName    string
		fileContent string
		expected    []VariableInfo
		wantErr     bool
	}{
		{
			name:        "simple field",
			fileName:    "test.tmpl",
			fileContent: `Hello {{.Name}}!`,
			expected:    []VariableInfo{{Name: "Name", DefaultValue: ""}},
			wantErr:     false,
		},
		{
			name:        "custom function without default",
			fileName:    "test.tmpl",
			fileContent: `{{getv "username"}}`,
			expected:    []VariableInfo{{Name: "username", DefaultValue: ""}},
			wantErr:     false,
		},
		{
			name:        "custom function with default",
			fileName:    "test.tmpl",
			fileContent: `{{getv "username" "default_user"}}`,
			expected:    []VariableInfo{{Name: "username", DefaultValue: "default_user"}},
			wantErr:     false,
		},
		{
			name:        "mixed variables",
			fileName:    "test.tmpl",
			fileContent: `Hello {{.Name}}, {{getv "email" "default@example.com"}}`,
			expected: []VariableInfo{
				{Name: "Name", DefaultValue: ""},
				{Name: "email", DefaultValue: "default@example.com"},
			},
			wantErr: false,
		},
		{
			name:        "multiple custom functions with defaults",
			fileName:    "test.tmpl",
			fileContent: `{{getv "host" "localhost"}}:{{getv "port" "8080"}}`,
			expected: []VariableInfo{
				{Name: "host", DefaultValue: "localhost"},
				{Name: "port", DefaultValue: "8080"},
			},
			wantErr: false,
		},
		{
			name:        "non-matching function should process arguments",
			fileName:    "test.tmpl",
			fileContent: `{{printf "%s" "hello"}}`,
			expected:    []VariableInfo{},
			wantErr:     false,
		},
		{
			name:        "nested with custom functions",
			fileName:    "test.tmpl",
			fileContent: `{{if .Enabled}}{{getv "name" "guest"}}{{else}}{{getv "name" "anonymous"}}{{end}}`,
			expected: []VariableInfo{
				{Name: "Enabled", DefaultValue: ""},
				{Name: "name", DefaultValue: "guest"},
				{Name: "name", DefaultValue: "anonymous"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ExtractVariablesWithDefaults(tt.fileName, tt.fileContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractVariablesWithDefaults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Handle empty slice comparison
				if len(got) == 0 && len(tt.expected) == 0 {
					// Both empty, this is fine
				} else if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("ExtractVariablesWithDefaults() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

func TestParser_createMinimalFuncMap_Pure(t *testing.T) {
	matcher := &DefaultFunctionMatcher{}
	parser := NewParser(matcher)

	funcMap := parser.createMinimalFuncMap()

	if funcMap == nil {
		t.Fatal("createMinimalFuncMap() returned nil")
	}

	// Test that all expected functions exist
	expectedFuncs := []string{"getv", "exists", "get", "jsonv"}
	for _, funcName := range expectedFuncs {
		if _, exists := funcMap[funcName]; !exists {
			t.Errorf("createMinimalFuncMap() missing function %s", funcName)
		}
	}

	// Test that functions work without errors (they should return empty/default values)
	if result := funcMap["getv"].(func(string, ...string) string)("test", "default"); result != "" {
		t.Errorf("Expected minimal getv to return empty string, got %v", result)
	}

	if result := funcMap["exists"].(func(string) bool)("test"); result != false {
		t.Errorf("Expected minimal exists to return false, got %v", result)
	}
}

func TestVariableInfo_Structure(t *testing.T) {
	varInfo := VariableInfo{
		Name:         "test",
		DefaultValue: "default",
	}

	if varInfo.Name != "test" {
		t.Errorf("VariableInfo.Name = %v, want %v", varInfo.Name, "test")
	}

	if varInfo.DefaultValue != "default" {
		t.Errorf("VariableInfo.DefaultValue = %v, want %v", varInfo.DefaultValue, "default")
	}
}

func TestParser_FunctionMatcherInterface_Pure(t *testing.T) {
	// Test that we can create a custom function matcher
	customMatcher := struct {
		FunctionMatcher
	}{
		// Empty struct implementing the interface
	}

	// This should compile without error
	parser := NewParser(customMatcher)
	if parser == nil {
		t.Fatal("NewParser() with custom matcher returned nil")
	}

	if parser.functionMatcher != customMatcher {
		t.Errorf("NewParser() functionMatcher = %v, want %v", parser.functionMatcher, customMatcher)
	}
}

func TestParser_ComplexTemplate_Pure(t *testing.T) {
	matcher := &DefaultFunctionMatcher{}
	parser := NewParser(matcher)

	// Test a complex template with multiple constructs
	complexTemplate := `
<!DOCTYPE html>
<html>
<head>
    <title>{{getv "title" "Default Title"}}</title>
</head>
<body>
    {{if .User}}
        <h1>Welcome {{.User.Name}}</h1>
        <p>Email: {{.User.Email}}</p>
    {{else}}
        <h1>{{getv "welcome" "Welcome Guest"}}</h1>
    {{end}}

    <ul>
    {{range .Items}}
        <li>{{.Name}} - {{.Description}}</li>
    {{end}}
    </ul>

    {{with .Settings}}
        <p>Theme: {{getv .Theme "light"}}</p>
        <p>Language: {{getv .Language "en"}}</p>
    {{end}}
</body>
</html>`

	expectedVars := []string{
		"title", "User", "welcome", "Items", "Name", "Description", "Settings", "Theme", "Language",
	}

	got, err := parser.ExtractVariables("complex.tmpl", complexTemplate)
	if err != nil {
		t.Fatalf("ExtractVariables() error = %v", err)
	}

	// Check that we got at least the expected variables
	for _, expected := range expectedVars {
		found := false
		for _, gotVar := range got {
			if gotVar == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected variable %s not found in result %v", expected, got)
		}
	}
}
