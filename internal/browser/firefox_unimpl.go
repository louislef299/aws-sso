//go:build !darwin && !linux

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
