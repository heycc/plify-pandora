//go:build !official
// +build !official

package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"path"
	"strconv"
	"strings"
	"text/template/parse"
	"time"
)

func init() {
	// Register Confd-style functions on initialization
	// This only happens when the "official" build tag is NOT present
	registerConfdFunctions()
}

// registerConfdFunctions registers all Confd-style template functions
func registerConfdFunctions() {
	registry := GetGlobalRegistry()

	// base - Base function (path.Base) - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "base",
		Description:           "Returns the last element of path",
		Handler:               baseMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// split - Split function (strings.Split) - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "split",
		Description:           "Splits a string into substrings separated by separator",
		Handler:               splitMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// json - Parse JSON object - already handled in functions_custom.go
	// This is just a duplicate registration, but we'll keep it for completeness
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "json",
		Description:           "Parse JSON variable and return as map",
		Handler:               jsonConfdMinimalHandler,
		Extractor:             extractSingleStringArgVariable,
		ExtractorWithDefaults: extractSingleStringArgVariableInfo,
	})

	// jsonArray - Parse JSON array - already handled in functions_custom.go
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "jsonArray",
		Description:           "Parse JSON variable and return as array",
		Handler:               jsonArrayConfdMinimalHandler,
		Extractor:             extractSingleStringArgVariable,
		ExtractorWithDefaults: extractSingleStringArgVariableInfo,
	})

	// dir - Directory function (path.Dir) - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "dir",
		Description:           "Returns all but the last element of path",
		Handler:               dirMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// map - Create map from key-value pairs - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "map",
		Description:           "Create a map from key-value pairs",
		Handler:               mapMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// join - Join function (strings.Join) - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "join",
		Description:           "Joins array elements into a string with separator",
		Handler:               joinMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// datetime - Current time - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "datetime",
		Description:           "Returns current time",
		Handler:               datetimeMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// toUpper - Convert to uppercase - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "toUpper",
		Description:           "Converts string to uppercase",
		Handler:               toUpperMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// toLower - Convert to lowercase - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "toLower",
		Description:           "Converts string to lowercase",
		Handler:               toLowerMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// replace - Replace substrings - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "replace",
		Description:           "Replaces old string with new string",
		Handler:               replaceMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// contains - Check if string contains substring - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "contains",
		Description:           "Checks if string contains substring",
		Handler:               containsMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// base64Encode - Base64 encode - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "base64Encode",
		Description:           "Base64 encodes a string",
		Handler:               base64EncodeMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// base64Decode - Base64 decode - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "base64Decode",
		Description:           "Base64 decodes a string",
		Handler:               base64DecodeMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// trimSuffix - Trim suffix - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "trimSuffix",
		Description:           "Trims suffix from string",
		Handler:               trimSuffixMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// parseBool - Parse boolean - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "parseBool",
		Description:           "Parses string to boolean",
		Handler:               parseBoolMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// reverse - Reverse array - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "reverse",
		Description:           "Reverses an array",
		Handler:               reverseMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// add - Add numbers - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "add",
		Description:           "Adds two numbers",
		Handler:               addMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// sub - Subtract numbers - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "sub",
		Description:           "Subtracts two numbers",
		Handler:               subMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// div - Divide numbers - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "div",
		Description:           "Divides two numbers",
		Handler:               divMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// mod - Modulo operation - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "mod",
		Description:           "Modulo operation on two numbers",
		Handler:               modMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// mul - Multiply numbers - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "mul",
		Description:           "Multiplies two numbers",
		Handler:               mulMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// seq - Generate sequence - extracts variables from first argument
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "seq",
		Description:           "Generates sequence of integers",
		Handler:               seqMinimalHandler,
		Extractor:             extractFirstArgVariable,
		ExtractorWithDefaults: extractFirstArgVariableInfo,
	})

	// atoi - String to integer - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "atoi",
		Description:           "Converts string to integer",
		Handler:               atoiMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})
}

// Minimal handlers for parsing (don't need actual variable values)
func baseMinimalHandler(s string) string {
	return ""
}

func splitMinimalHandler(s, sep string) []string {
	return []string{}
}

func jsonConfdMinimalHandler(data string) (map[string]interface{}, error) {
	return nil, nil
}

func jsonArrayConfdMinimalHandler(data string) ([]interface{}, error) {
	return nil, nil
}

func dirMinimalHandler(s string) string {
	return ""
}

func mapMinimalHandler(values ...interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func joinMinimalHandler(elems []string, sep string) string {
	return ""
}

func datetimeMinimalHandler() time.Time {
	return time.Time{}
}

func toUpperMinimalHandler(s string) string {
	return ""
}

func toLowerMinimalHandler(s string) string {
	return ""
}

func replaceMinimalHandler(s, old, new string, n int) string {
	return ""
}

func containsMinimalHandler(s, substr string) bool {
	return false
}

func base64EncodeMinimalHandler(data string) string {
	return ""
}

func base64DecodeMinimalHandler(data string) (string, error) {
	return "", nil
}

func trimSuffixMinimalHandler(s, suffix string) string {
	return ""
}

func parseBoolMinimalHandler(str string) (bool, error) {
	return false, nil
}

func reverseMinimalHandler(values interface{}) interface{} {
	return nil
}

func addMinimalHandler(a, b int) int {
	return 0
}

func subMinimalHandler(a, b int) int {
	return 0
}

func divMinimalHandler(a, b int) int {
	return 0
}

func modMinimalHandler(a, b int) int {
	return 0
}

func mulMinimalHandler(a, b int) int {
	return 0
}

func seqMinimalHandler(first, last int) []int {
	return []int{}
}

func atoiMinimalHandler(s string) (int, error) {
	return 0, nil
}

// Actual handlers for rendering (use variable values)
func baseRenderHandler(_ map[string]interface{}) func(s string) string {
	return func(s string) string {
		return path.Base(s)
	}
}

func splitRenderHandler(_ map[string]interface{}) func(s, sep string) []string {
	return func(s, sep string) []string {
		return strings.Split(s, sep)
	}
}

func jsonConfdRenderHandler(variables map[string]interface{}) func(key string) (map[string]interface{}, error) {
	return func(key string) (map[string]interface{}, error) {
		if val, exists := variables[key]; exists {
			if strVal, ok := val.(string); ok {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(strVal), &result)
				return result, err
			}
		}
		return nil, errors.New("key not found")
	}
}

func jsonArrayConfdRenderHandler(variables map[string]interface{}) func(key string) ([]interface{}, error) {
	return func(key string) ([]interface{}, error) {
		if val, exists := variables[key]; exists {
			if strVal, ok := val.(string); ok {
				var result []interface{}
				err := json.Unmarshal([]byte(strVal), &result)
				return result, err
			}
		}
		return nil, errors.New("key not found")
	}
}

func dirRenderHandler(_ map[string]interface{}) func(s string) string {
	return func(s string) string {
		return path.Dir(s)
	}
}

func mapRenderHandler(_ map[string]interface{}) func(values ...interface{}) (map[string]interface{}, error) {
	return func(values ...interface{}) (map[string]interface{}, error) {
		if len(values)%2 != 0 {
			return nil, errors.New("invalid map call")
		}
		dict := make(map[string]interface{}, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			key, ok := values[i].(string)
			if !ok {
				return nil, errors.New("map keys must be strings")
			}
			dict[key] = values[i+1]
		}
		return dict, nil
	}
}

func joinRenderHandler(_ map[string]interface{}) func(elems []string, sep string) string {
	return func(elems []string, sep string) string {
		return strings.Join(elems, sep)
	}
}

func datetimeRenderHandler(_ map[string]interface{}) func() time.Time {
	return func() time.Time {
		return time.Now()
	}
}

func toUpperRenderHandler(_ map[string]interface{}) func(s string) string {
	return func(s string) string {
		return strings.ToUpper(s)
	}
}

func toLowerRenderHandler(_ map[string]interface{}) func(s string) string {
	return func(s string) string {
		return strings.ToLower(s)
	}
}

func replaceRenderHandler(_ map[string]interface{}) func(s, old, new string, n int) string {
	return func(s, old, new string, n int) string {
		return strings.Replace(s, old, new, n)
	}
}

func containsRenderHandler(_ map[string]interface{}) func(s, substr string) bool {
	return func(s, substr string) bool {
		return strings.Contains(s, substr)
	}
}

func base64EncodeRenderHandler(_ map[string]interface{}) func(data string) string {
	return func(data string) string {
		return base64.StdEncoding.EncodeToString([]byte(data))
	}
}

func base64DecodeRenderHandler(_ map[string]interface{}) func(data string) (string, error) {
	return func(data string) (string, error) {
		s, err := base64.StdEncoding.DecodeString(data)
		return string(s), err
	}
}

func trimSuffixRenderHandler(_ map[string]interface{}) func(s, suffix string) string {
	return func(s, suffix string) string {
		return strings.TrimSuffix(s, suffix)
	}
}

func parseBoolRenderHandler(_ map[string]interface{}) func(str string) (bool, error) {
	return func(str string) (bool, error) {
		return strconv.ParseBool(str)
	}
}

func reverseRenderHandler(_ map[string]interface{}) func(values interface{}) interface{} {
	return func(values interface{}) interface{} {
		switch v := values.(type) {
		case []string:
			for left, right := 0, len(v)-1; left < right; left, right = left+1, right-1 {
				v[left], v[right] = v[right], v[left]
			}
			return v
		case []interface{}:
			for left, right := 0, len(v)-1; left < right; left, right = left+1, right-1 {
				v[left], v[right] = v[right], v[left]
			}
			return v
		}
		return values
	}
}

func addRenderHandler(_ map[string]interface{}) func(a, b int) int {
	return func(a, b int) int {
		return a + b
	}
}

func subRenderHandler(_ map[string]interface{}) func(a, b int) int {
	return func(a, b int) int {
		return a - b
	}
}

func divRenderHandler(_ map[string]interface{}) func(a, b int) int {
	return func(a, b int) int {
		return a / b
	}
}

func modRenderHandler(_ map[string]interface{}) func(a, b int) int {
	return func(a, b int) int {
		return a % b
	}
}

func mulRenderHandler(_ map[string]interface{}) func(a, b int) int {
	return func(a, b int) int {
		return a * b
	}
}

func seqRenderHandler(_ map[string]interface{}) func(first, last int) []int {
	return func(first, last int) []int {
		var arr []int
		for i := first; i <= last; i++ {
			arr = append(arr, i)
		}
		return arr
	}
}

func atoiRenderHandler(_ map[string]interface{}) func(s string) (int, error) {
	return func(s string) (int, error) {
		return strconv.Atoi(s)
	}
}

// Variable extractors
func extractNoVariables(args []parse.Node, cycle int) ([]string, error) {
	return []string{}, nil
}

func extractNoVariablesInfo(args []parse.Node, cycle int) ([]VariableInfo, error) {
	return []VariableInfo{}, nil
}

// Only json and jsonArray need variable extraction since they access variables
func extractSingleStringArgVariable(args []parse.Node, cycle int) ([]string, error) {
	return extractStringArgVariable(args, cycle, 1)
}

func extractSingleStringArgVariableInfo(args []parse.Node, cycle int) ([]VariableInfo, error) {
	return extractStringArgVariableWithDefaults(args, cycle, 1, -1)
}

// Extract variables from first argument (for transformation functions like base, toUpper, etc.)
// These functions operate on their first argument, which may be a variable reference
// Unlike extractStringArgVariable, this does NOT treat string literals as variable names
func extractFirstArgVariable(args []parse.Node, cycle int) ([]string, error) {
	return extractArgVariable(args, cycle, 1, false)
}

func extractFirstArgVariableInfo(args []parse.Node, cycle int) ([]VariableInfo, error) {
	return extractArgVariableWithDefaults(args, cycle, 1, -1, false)
}

// CreateConfdRenderFuncMap creates function map with actual variable values for rendering Confd functions
func CreateConfdRenderFuncMap(variables map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"base":         baseRenderHandler(variables),
		"split":        splitRenderHandler(variables),
		"json":         jsonConfdRenderHandler(variables),
		"jsonArray":    jsonArrayConfdRenderHandler(variables),
		"dir":          dirRenderHandler(variables),
		"map":          mapRenderHandler(variables),
		"join":         joinRenderHandler(variables),
		"datetime":     datetimeRenderHandler(variables),
		"toUpper":      toUpperRenderHandler(variables),
		"toLower":      toLowerRenderHandler(variables),
		"replace":      replaceRenderHandler(variables),
		"contains":     containsRenderHandler(variables),
		"base64Encode": base64EncodeRenderHandler(variables),
		"base64Decode": base64DecodeRenderHandler(variables),
		"trimSuffix":   trimSuffixRenderHandler(variables),
		"parseBool":    parseBoolRenderHandler(variables),
		"reverse":      reverseRenderHandler(variables),
		"add":          addRenderHandler(variables),
		"sub":          subRenderHandler(variables),
		"div":          divRenderHandler(variables),
		"mod":          modRenderHandler(variables),
		"mul":          mulRenderHandler(variables),
		"seq":          seqRenderHandler(variables),
		"atoi":         atoiRenderHandler(variables),
	}
}
