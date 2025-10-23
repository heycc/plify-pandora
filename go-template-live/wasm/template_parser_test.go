//go:build !js
// +build !js

package main

import (
	"reflect"
	"testing"
)

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
		wantErr          bool
	}{
		{
			name:             "simple field",
			fileName:         "test.tmpl",
			fileContent:      `Hello {{.Name}}!`,
			expectedCustom:   []string{"Name"},
			expectedOfficial: []string{"Name"},
			wantErr:          false,
		},
		{
			name:             "nested field",
			fileName:         "test.tmpl",
			fileContent:      `Hello {{.User.Name}}!`,
			expectedCustom:   []string{"User.Name"},
			expectedOfficial: []string{"User.Name"},
			wantErr:          false,
		},
		{
			name:             "custom function with string parameter",
			fileName:         "test.tmpl",
			fileContent:      `{{getv "username"}}`,
			expectedCustom:   []string{"username"},
			expectedOfficial: []string{}, // Official parser ignores custom functions
			wantErr:          false,
		},
		{
			name:             "custom function with field parameter",
			fileName:         "test.tmpl",
			fileContent:      `{{getv .Key}}`,
			expectedCustom:   []string{"Key"},
			expectedOfficial: []string{"Key"}, // Field parameters still work
			wantErr:          false,
		},
		{
			name:             "multiple variables",
			fileName:         "test.tmpl",
			fileContent:      `Hello {{.Name}}, your email is {{.Email}} and username is {{getv "username"}}`,
			expectedCustom:   []string{"Name", "Email", "username"},
			expectedOfficial: []string{"Name", "Email"}, // Ignores custom function variables
			wantErr:          false,
		},
		{
			name:             "if statement",
			fileName:         "test.tmpl",
			fileContent:      `{{if .Enabled}}{{.Name}}{{end}}`,
			expectedCustom:   []string{"Enabled", "Name"},
			expectedOfficial: []string{"Enabled", "Name"},
			wantErr:          false,
		},
		{
			name:             "range statement",
			fileName:         "test.tmpl",
			fileContent:      `{{range .Items}}{{.Name}}{{end}}`,
			expectedCustom:   []string{"Items", "Name"},
			expectedOfficial: []string{"Items", "Name"},
			wantErr:          false,
		},
		{
			name:             "with statement",
			fileName:         "test.tmpl",
			fileContent:      `{{with .User}}{{.Name}}{{end}}`,
			expectedCustom:   []string{"User", "Name"},
			expectedOfficial: []string{"User", "Name"},
			wantErr:          false,
		},
		{
			name:             "non-matching function should process arguments",
			fileName:         "test.tmpl",
			fileContent:      `{{printf "%s" "hello"}}`,
			expectedCustom:   []string{},
			expectedOfficial: []string{},
			wantErr:          false,
		},
		{
			name:             "invalid template",
			fileName:         "test.tmpl",
			fileContent:      `{{.unclosed`,
			expectedCustom:   nil,
			expectedOfficial: nil,
			wantErr:          true,
		},
		{
			name:             "template with else branch",
			fileName:         "test.tmpl",
			fileContent:      `{{if .Active}}{{.Name}}{{else}}{{.DefaultName}}{{end}}`,
			expectedCustom:   []string{"Active", "Name", "DefaultName"},
			expectedOfficial: []string{"Active", "Name", "DefaultName"},
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_custom", func(t *testing.T) {
			got, err := parserCustom.ExtractVariables(tt.fileName, tt.fileContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractVariables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
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
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractVariables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
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