//go:build windows
// +build windows

package browser

import (
	"context"
)

const macosFirefoxPath = "/Applications/Firefox.app/Contents/MacOS/firefox"

type Firefox struct {
	private bool
}

func (f *Firefox) OpenURL(ctx context.Context, url string) error {
	if f.private {
		return open(ctx, macosFirefoxPath, "--private-window", url)
	}
	return open(ctx, macosFirefoxPath, url)
}
