//go:build linux
// +build linux

package browser

import (
	"context"

	los "github.com/louislef299/aws-sso/pkg/v1/os"
)

var linuxFirefoxPaths = []string{`/usr/bin/firefox`, `/../../mnt/c/Program Files/Mozilla Firefox/firefox.exe`}

type Firefox struct {
	private bool
}

func (f *Firefox) OpenURL(ctx context.Context, url string) error {
	var linuxPath string
	for _, path := range linuxFirefoxPaths {
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

func (f *Firefox) Type() string { return "firefox" }
