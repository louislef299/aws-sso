//go:build darwin
// +build darwin

package browser

import (
	"context"
	"os/exec"
)

func open(ctx context.Context, filePath string, cmds ...string) error {
	path, err := exec.LookPath(filePath)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, path, cmds...)
	return cmd.Run()
}
