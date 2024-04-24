//go:build darwin
// +build darwin

package browser

import (
	"context"
)

const macosBravePath = "/Applications/Brave Browser.app/Contents/MacOS/Brave Browser"

type Brave struct {
	private bool
}

func (f *Brave) OpenURL(ctx context.Context, url string) error {
	if f.private {
		return open(ctx, macosBravePath, "--incognito", url)
	}
	return open(ctx, macosBravePath, url)
}
