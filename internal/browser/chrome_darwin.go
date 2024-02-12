//go:build darwin
// +build darwin

package browser

import (
	"context"
)

const macosPath = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"

type Chrome struct {
	private bool
}

func (f *Chrome) OpenURL(ctx context.Context, url string) error {
	if f.private {
		return open(ctx, macosPath, "-incognito", url)
	}
	return open(ctx, macosPath, url)
}

func (f *Chrome) Type() string { return "chrome" }
