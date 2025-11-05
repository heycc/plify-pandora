//go:build js && !custom && !confd && !official
// +build js,!custom,!confd,!official

package main

// CreateRenderFuncMap provides a no-op function map for builds that do not
// include any of the custom build tags. This ensures js/wasm builds without
// additional tags still compile and satisfy references from wasm_handlers.go.
func CreateRenderFuncMap(variables map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{}
}
