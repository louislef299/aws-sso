//go:build !darwin
// +build !darwin

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
