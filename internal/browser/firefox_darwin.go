//go:build darwin
// +build darwin

package browser

import (
	"context"
)

const macosPath = "/Applications/Firefox.app/Contents/MacOS/firefox"

type Firefox struct {
	private bool
}

func (f *Firefox) OpenURL(ctx context.Context, url string) error {
	if f.private {
		return open(ctx, macosPath, "--private-window", url)
	}
	return open(ctx, macosPath, url)
}

func (f *Firefox) Type() string { return "firefox" }
