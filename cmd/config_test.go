package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfigCommands(t *testing.T) {
	// Setup: use a temp config file
	tmpDir := t.TempDir()
	tmpConfig := filepath.Join(tmpDir, "test.toml")
	viper.SetConfigFile(tmpConfig)
	viper.SetConfigType("toml")
	err := viper.WriteConfigAs(tmpConfig)
	assert.NoError(t, err)

	// Test 'set' command
	setCmd := rootCmd
	buf := new(bytes.Buffer)
	setCmd.SetOut(buf)
	setCmd.SetArgs([]string{"config", "set", "core.testKey", "testValue"})
	err = setCmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, viper.GetString("core.testKey"), "testValue")

	// Test 'get' command
	getCmd := rootCmd
	buf.Reset()
	getCmd.SetOut(buf)
	getCmd.SetArgs([]string{"config", "get", "core.testKey"})
	err = getCmd.Execute()
	assert.NoError(t, err)
	t.Log(buf.String())
	assert.Contains(t, buf.String(), "testValue")

	// Test 'unset' command
	unsetCmd := rootCmd
	buf.Reset()
	unsetCmd.SetOut(buf)
	unsetCmd.SetArgs([]string{"config", "unset", "core.testKey"})
	err = unsetCmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "", viper.GetString("core.testKey"))

	// Test 'list' command
	listCmd := rootCmd
	buf.Reset()
	listCmd.SetOut(buf)
	listCmd.SetArgs([]string{"config", "ls"})
	err = listCmd.Execute()
	assert.NoError(t, err)

	// Cleanup
	os.Remove(tmpConfig)
}

func TestParsePlugins(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []string
	}{
		// String inputs
		{
			name:     "comma-separated with brackets",
			input:    "[oidc,eks]",
			expected: []string{"oidc", "eks"},
		},
		{
			name:     "comma-separated without brackets",
			input:    "oidc,eks",
			expected: []string{"oidc", "eks"},
		},
		{
			name:     "comma-separated with spaces",
			input:    "oidc, eks, ecr",
			expected: []string{"oidc", "eks", "ecr"},
		},
		{
			name:     "comma-separated with quotes",
			input:    "\"oidc\",\"eks\",\"ecr\"",
			expected: []string{"oidc", "eks", "ecr"},
		},
		{
			name:     "comma-separated with single quotes",
			input:    "'oidc','eks','ecr'",
			expected: []string{"oidc", "eks", "ecr"},
		},
		{
			name:     "space-separated",
			input:    "oidc eks ecr",
			expected: []string{"oidc", "eks", "ecr"},
		},
		{
			name:     "single plugin",
			input:    "oidc",
			expected: []string{"oidc"},
		},
		{
			name:     "single plugin with brackets",
			input:    "[oidc]",
			expected: []string{"oidc"},
		},
		{
			name:     "mixed brackets and quotes",
			input:    "[\"oidc\",\"eks\"]",
			expected: []string{"oidc", "eks"},
		},
		{
			name:     "extra whitespace",
			input:    "  oidc  ,  eks  ",
			expected: []string{"oidc", "eks"},
		},
		{
			name:     "all three plugins",
			input:    "oidc,eks,ecr",
			expected: []string{"oidc", "eks", "ecr"},
		},
		// Array inputs (TOML arrays)
		{
			name:     "proper TOML array",
			input:    []string{"oidc", "eks", "ecr"},
			expected: []string{"oidc", "eks", "ecr"},
		},
		{
			name:     "single element array",
			input:    []string{"oidc"},
			expected: []string{"oidc"},
		},
		{
			name:     "array with concatenated string",
			input:    []string{"oidc,eks"},
			expected: []string{"oidc", "eks"},
		},
		{
			name:     "array with space-separated string",
			input:    []string{"oidc eks"},
			expected: []string{"oidc", "eks"},
		},
		// Interface slice inputs (how viper reads TOML arrays)
		{
			name:     "interface slice",
			input:    []interface{}{"oidc", "eks", "ecr"},
			expected: []string{"oidc", "eks", "ecr"},
		},
		{
			name:     "single element interface slice",
			input:    []interface{}{"oidc"},
			expected: []string{"oidc"},
		},
		{
			name:     "empty interface slice",
			input:    []interface{}{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePlugins(tt.input)
			assert.Equal(t, tt.expected, result, "parsePlugins(%v) failed", tt.input)
		})
	}
}

func TestParsePlugins_EmptyInput(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{name: "nil", input: nil},
		{name: "empty string", input: ""},
		{name: "only whitespace", input: "   "},
		{name: "empty brackets", input: "[]"},
		{name: "brackets with whitespace", input: "[  ]"},
		{name: "empty array", input: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePlugins(tt.input)
			assert.Empty(t, result, "parsePlugins(%v) should return empty slice", tt.input)
		})
	}
}

func TestParsePluginString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "comma-separated",
			input:    "oidc,eks,ecr",
			expected: []string{"oidc", "eks", "ecr"},
		},
		{
			name:     "space-separated",
			input:    "oidc eks ecr",
			expected: []string{"oidc", "eks", "ecr"},
		},
		{
			name:     "with brackets",
			input:    "[oidc,eks]",
			expected: []string{"oidc", "eks"},
		},
		{
			name:     "with quotes",
			input:    `"oidc","eks"`,
			expected: []string{"oidc", "eks"},
		},
		{
			name:     "single plugin",
			input:    "oidc",
			expected: []string{"oidc"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePluginString(tt.input)
			assert.Equal(t, tt.expected, result, "parsePluginString(%q) failed", tt.input)
		})
	}
}

func TestCleanPluginName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "plain name", input: "oidc", expected: "oidc"},
		{name: "with whitespace", input: "  oidc  ", expected: "oidc"},
		{name: "with double quotes", input: `"oidc"`, expected: "oidc"},
		{name: "with single quotes", input: "'oidc'", expected: "oidc"},
		{name: "with mixed whitespace and quotes", input: `  "oidc"  `, expected: "oidc"},
		{name: "empty string", input: "", expected: ""},
		{name: "only quotes", input: `""`, expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanPluginName(tt.input)
			assert.Equal(t, tt.expected, result, "cleanPluginName(%q) failed", tt.input)
		})
	}
}

func TestValidateInput_CorePlugins(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "valid comma-separated",
			input:    "oidc,eks",
			expected: []string{"oidc", "eks"},
		},
		{
			name:     "valid space-separated",
			input:    "oidc eks",
			expected: []string{"oidc", "eks"},
		},
		{
			name:     "valid with brackets",
			input:    "[oidc,eks]",
			expected: []string{"oidc", "eks"},
		},
		{
			name:     "all plugins",
			input:    "oidc,eks,ecr",
			expected: []string{"oidc", "eks", "ecr"},
		},
		{
			name:     "single valid plugin",
			input:    "oidc",
			expected: []string{"oidc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup: use a temp config file
			tmpDir := t.TempDir()
			tmpConfig := filepath.Join(tmpDir, "test.toml")
			viper.SetConfigFile(tmpConfig)
			viper.SetConfigType("toml")
			err := viper.WriteConfigAs(tmpConfig)
			assert.NoError(t, err)

			// Clear any previous value
			viper.Set("core.plugins", nil)

			// Run validateInput
			validateInput("core.plugins", tt.input)

			// Check that the value was set correctly as a slice
			result := viper.GetStringSlice("core.plugins")
			assert.Equal(t, tt.expected, result, "validateInput should set correct plugin list")

			// Cleanup
			os.Remove(tmpConfig)
		})
	}
}

func TestRegisteredPluginDrivers(t *testing.T) {
	// Test that registeredPluginDrivers returns a non-empty list
	drivers := registeredPluginDrivers()
	assert.NotEmpty(t, drivers, "registeredPluginDrivers should return at least one driver")

	// Test that it includes the expected plugins
	assert.Contains(t, drivers, "oidc", "Should include oidc driver")
	assert.Contains(t, drivers, "eks", "Should include eks driver")
	assert.Contains(t, drivers, "ecr", "Should include ecr driver")
}
