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
	"time"
	"text/template/parse"
)

func init() {
	// Register Confd-style functions on initialization
	// This only happens when the "official" build tag is NOT present
	registerConfdFunctions()
}

// registerConfdFunctions registers all Confd-style template functions
func registerConfdFunctions() {
	registry := GetGlobalRegistry()

	// base - Base function (path.Base) - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "base",
		Description:           "Returns the last element of path",
		Handler:               baseMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// split - Split function (strings.Split) - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "split",
		Description:           "Splits a string into substrings separated by separator",
		Handler:               splitMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
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

	// dir - Directory function (path.Dir) - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "dir",
		Description:           "Returns all but the last element of path",
		Handler:               dirMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
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

	// toUpper - Convert to uppercase - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "toUpper",
		Description:           "Converts string to uppercase",
		Handler:               toUpperMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// toLower - Convert to lowercase - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "toLower",
		Description:           "Converts string to lowercase",
		Handler:               toLowerMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// replace - Replace substrings - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "replace",
		Description:           "Replaces old string with new string",
		Handler:               replaceMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// contains - Check if string contains substring - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "contains",
		Description:           "Checks if string contains substring",
		Handler:               containsMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// base64Encode - Base64 encode - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "base64Encode",
		Description:           "Base64 encodes a string",
		Handler:               base64EncodeMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// base64Decode - Base64 decode - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "base64Decode",
		Description:           "Base64 decodes a string",
		Handler:               base64DecodeMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// trimSuffix - Trim suffix - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "trimSuffix",
		Description:           "Trims suffix from string",
		Handler:               trimSuffixMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// parseBool - Parse boolean - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "parseBool",
		Description:           "Parses string to boolean",
		Handler:               parseBoolMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// reverse - Reverse array - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "reverse",
		Description:           "Reverses an array",
		Handler:               reverseMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// add - Add numbers - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "add",
		Description:           "Adds two numbers",
		Handler:               addMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// sub - Subtract numbers - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "sub",
		Description:           "Subtracts two numbers",
		Handler:               subMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// div - Divide numbers - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "div",
		Description:           "Divides two numbers",
		Handler:               divMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// mod - Modulo operation - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "mod",
		Description:           "Modulo operation on two numbers",
		Handler:               modMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// mul - Multiply numbers - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "mul",
		Description:           "Multiplies two numbers",
		Handler:               mulMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
	})

	// seq - Generate sequence - pure utility, no variable extraction
	registry.RegisterFunction(&FunctionDefinition{
		Name:                  "seq",
		Description:           "Generates sequence of integers",
		Handler:               seqMinimalHandler,
		Extractor:             extractNoVariables,
		ExtractorWithDefaults: extractNoVariablesInfo,
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

// CreateConfdRenderFuncMap creates function map with actual variable values for rendering Confd functions
func CreateConfdRenderFuncMap(variables map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"base":        baseRenderHandler(variables),
		"split":       splitRenderHandler(variables),
		"json":        jsonConfdRenderHandler(variables),
		"jsonArray":   jsonArrayConfdRenderHandler(variables),
		"dir":         dirRenderHandler(variables),
		"map":         mapRenderHandler(variables),
		"join":        joinRenderHandler(variables),
		"datetime":    datetimeRenderHandler(variables),
		"toUpper":     toUpperRenderHandler(variables),
		"toLower":     toLowerRenderHandler(variables),
		"replace":     replaceRenderHandler(variables),
		"contains":    containsRenderHandler(variables),
		"base64Encode": base64EncodeRenderHandler(variables),
		"base64Decode": base64DecodeRenderHandler(variables),
		"trimSuffix":  trimSuffixRenderHandler(variables),
		"parseBool":   parseBoolRenderHandler(variables),
		"reverse":     reverseRenderHandler(variables),
		"add":         addRenderHandler(variables),
		"sub":         subRenderHandler(variables),
		"div":         divRenderHandler(variables),
		"mod":         modRenderHandler(variables),
		"mul":         mulRenderHandler(variables),
		"seq":         seqRenderHandler(variables),
		"atoi":        atoiRenderHandler(variables),
	}
}