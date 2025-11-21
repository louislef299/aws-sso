package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/louislef299/knot/internal/envs"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRootCmdConfig(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	path := filepath.Join("testdata", "config.toml")
	rootCmd.SetArgs([]string{"--config", path})

	// assert file exists first
	_, err := os.ReadFile(path)
	assert.NoError(t, err)

	err = rootCmd.Execute()
	assert.NoError(t, err)
	assert.Equal(t, "knot", rootCmd.Use)

	// TODO: fix this test case
	//assert.Equal(t, "firefox", viper.GetString(envs.CORE_BROWSER))
}

func TestParsePluginsConfig(t *testing.T) {
	tests := []struct {
		name     string
		setValue interface{}
		expected []string
	}{
		{
			name:     "proper TOML array",
			setValue: []string{"oidc", "eks", "ecr"},
			expected: []string{"oidc", "eks", "ecr"},
		},
		{
			name:     "comma-separated string with brackets",
			setValue: "[oidc,eks]",
			expected: []string{"oidc", "eks"},
		},
		{
			name:     "comma-separated string without brackets",
			setValue: "oidc,eks",
			expected: []string{"oidc", "eks"},
		},
		{
			name:     "space-separated string",
			setValue: "oidc eks ecr",
			expected: []string{"oidc", "eks", "ecr"},
		},
		{
			name:     "comma-separated with spaces",
			setValue: "oidc, eks, ecr",
			expected: []string{"oidc", "eks", "ecr"},
		},
		{
			name:     "single plugin as string",
			setValue: "oidc",
			expected: []string{"oidc"},
		},
		{
			name:     "single plugin in array",
			setValue: []string{"oidc"},
			expected: []string{"oidc"},
		},
		{
			name:     "comma-separated with quotes",
			setValue: `"oidc","eks"`,
			expected: []string{"oidc", "eks"},
		},
		{
			name:     "space-separated with extra whitespace",
			setValue: "  oidc   eks  ",
			expected: []string{"oidc", "eks"},
		},
		{
			name:     "empty string",
			setValue: "",
			expected: []string{},
		},
		{
			name:     "empty array",
			setValue: []string{},
			expected: []string{},
		},
		{
			name:     "only whitespace",
			setValue: "   ",
			expected: []string{},
		},
		{
			name:     "empty brackets",
			setValue: "[]",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Temporarily set the value in viper
			oldValue := viper.Get(envs.CORE_PLUGINS)
			viper.Set(envs.CORE_PLUGINS, tt.setValue)
			defer viper.Set(envs.CORE_PLUGINS, oldValue)

			result := parsePluginsConfig()
			assert.Equal(t, tt.expected, result, "parsePluginsConfig() with value %v failed", tt.setValue)
		})
	}
}
