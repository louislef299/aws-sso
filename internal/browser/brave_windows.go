//go:build windows
// +build windows

package browser

import (
	"context"
)

const windowsBravePath = `\Program Files\BraveSoftware\Brave-Browser\Application\brave.exe`

type Brave struct {
	private bool
}

func (f *Brave) OpenURL(ctx context.Context, url string) error {
	if f.private {
		return open(ctx, windowsBravePath, "--incognito", url)
	}
	return open(ctx, windowsBravePath, url)
}
