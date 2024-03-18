//go:build !darwin && !linux && !windows
// +build !darwin,!linux,!windows

package browser

import (
	"context"
)

type Brave struct {
	private bool
}

func (f *Brave) OpenURL(ctx context.Context, url string) error {
	return ErrBrowserUnimplemented
}

func (f *Brave) Type() string { return "unimplemented" }
