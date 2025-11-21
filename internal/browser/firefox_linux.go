//go:build linux

package browser

import (
	"context"

	los "github.com/louislef299/knot/pkg/os"
)

var (
	linuxFirefoxPaths                 = []string{`/usr/bin/firefox`, `/../../mnt/c/Program Files/Mozilla Firefox/firefox.exe`}
	linuxFirefoxDeveloperEditionPaths = []string{`/usr/bin/firefox-dev`, `/../../mnt/c/Program Files/Mozilla Firefox/firefox-dev.exe`}
)

type Firefox struct {
	private   bool
	developer bool
}

func (f *Firefox) OpenURL(ctx context.Context, url string) error {
	paths := linuxFirefoxPaths
	if f.developer {
		paths = linuxFirefoxDeveloperEditionPaths
	}

	var linuxPath string
	for _, path := range paths {
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
