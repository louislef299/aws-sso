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
