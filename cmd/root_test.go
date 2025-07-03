package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

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
