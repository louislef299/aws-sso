//go:build !darwin
// +build !darwin

package browser

import (
	"context"
)

type Chrome struct {
	private bool
}

func (f *Chrome) OpenURL(ctx context.Context, url string) error {
	return ErrBrowserUnimplemented
}

func (f *Chrome) Type() string { return "unimplemented" }
