package main

// BuildConfig controls which functions are included in the build
// This allows us to build different WASM files with different function sets
type BuildConfig struct {
	IncludeCustomFunctions bool
}

// DefaultBuildConfig returns the default build configuration
// This includes all custom functions
func DefaultBuildConfig() *BuildConfig {
	return &BuildConfig{
		IncludeCustomFunctions: true,
	}
}

// OfficialBuildConfig returns build configuration for official functions only
func OfficialBuildConfig() *BuildConfig {
	return &BuildConfig{
		IncludeCustomFunctions: false,
	}
}