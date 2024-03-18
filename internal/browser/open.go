package browser

import (
	"context"
	"os"
	"os/exec"
)

func open(ctx context.Context, filePath string, cmds ...string) error {
	path, err := exec.LookPath(filePath)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, path, cmds...)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
