//go:build linux
// +build linux

package browser

import (
	"context"

	los "github.com/louislef299/aws-sso/pkg/v1/os"
)

var linuxBravePaths = []string{`/usr/bin/brave-browser`, `/../../mnt/c/Program Files/BraveSoftware/Brave-Browser/Application/brave.exe`}

type Brave struct {
	private bool
}

func (f *Brave) OpenURL(ctx context.Context, url string) error {
	var linuxPath string
	for _, path := range linuxBravePaths {
		if exists, err := los.IsFileOrFolderExisting(path); err == nil && exists {
			linuxPath = path
			break
		} else {
			continue
		}
	}
	if linuxPath == "" {
		return ErrCouldNotFindBrowser
	}

	if f.private {
		return open(ctx, linuxPath, "--incognito", url)
	}
	return open(ctx, linuxPath, url)
}
