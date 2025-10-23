//go:build js
// +build js

package main

import ()

// registerCallbacks registers the Go functions to be called from JavaScript
func registerCallbacks() {
	handler := NewWASMHandler()
	handler.RegisterCallbacks()
}

func main() {
	c := make(chan struct{}, 0)
	registerCallbacks()
	<-c
}
