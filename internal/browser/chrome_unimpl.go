//go:build !darwin && !linux && !windows

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
