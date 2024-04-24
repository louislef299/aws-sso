//go:build darwin
// +build darwin

package browser

import (
	"context"
)

const macosChromePath = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"

type Chrome struct {
	private bool
}

func (f *Chrome) OpenURL(ctx context.Context, url string) error {
	if f.private {
		return open(ctx, macosChromePath, "-incognito", url)
	}
	return open(ctx, macosChromePath, url)
}
