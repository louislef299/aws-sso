//go:build !darwin && !linux && !windows
// +build !darwin,!linux,!windows

package browser

import (
	"context"
)

type Firefox struct {
	private bool
}

func (f *Firefox) OpenURL(ctx context.Context, url string) error {
	return ErrBrowserUnimplemented
}

func (f *Firefox) Type() string { return "unimplemented" }
