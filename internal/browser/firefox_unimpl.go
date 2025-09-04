//go:build !darwin && !linux && !windows
// +build !darwin,!linux,!windows

package browser

import (
	"context"
)

type Firefox struct {
	private   bool
	developer bool
}

func (f *Firefox) OpenURL(ctx context.Context, url string) error {
	return ErrBrowserUnimplemented
}
