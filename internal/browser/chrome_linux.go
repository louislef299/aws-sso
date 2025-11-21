//go:build linux

package browser

import (
	"context"

	los "github.com/louislef299/knot/pkg/os"
)

var linuxChromePaths = []string{`/usr/bin/google-chrome`, `/../../mnt/c/Program Files/Google/Chrome/Application/chrome.exe`, `/../../mnt/c/Program Files (x86)/Google/Chrome/Application/chrome.exe`}

type Chrome struct {
	private bool
}

func (f *Chrome) OpenURL(ctx context.Context, url string) error {
	var linuxPath string
	for _, path := range linuxChromePaths {
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
