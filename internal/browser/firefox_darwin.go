//go:build darwin
// +build darwin

package browser

import (
	"context"
)

const (
	macosFirefoxPath                 = "/Applications/Firefox.app/Contents/MacOS/firefox"
	macosFirefoxDeveloperEditionPath = "/Applications/Firefox Developer Edition.app/Contents/MacOS/firefox"
)

type Firefox struct {
	private   bool
	developer bool
}

func (f *Firefox) OpenURL(ctx context.Context, url string) error {
	path := macosFirefoxPath
	if f.developer {
		path = macosFirefoxDeveloperEditionPath
	}

	if f.private {
		return open(ctx, path, "--private-window", url)
	}
	return open(ctx, path, url)
}
