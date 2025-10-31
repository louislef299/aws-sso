//go:build windows

package browser

import (
	"context"

	los "github.com/louislef299/aws-sso/pkg/os"
)

var windowsChromePaths = []string{`\Program Files\Google\Chrome\Application\chrome.exe`, `\Program Files (x86)\Google\Chrome\Application\chrome.exe`}

type Chrome struct {
	private bool
}

func (f *Chrome) OpenURL(ctx context.Context, url string) error {
	var windowsPath string
	for _, path := range windowsChromePaths {
		if exists, err := los.IsFileOrFolderExisting(path); err == nil && exists {
			windowsPath = path
			break
		} else {
			continue
		}
	}
	if windowsPath == "" {
		return ErrCouldNotFindBrowser
	}

	if f.private {
		return open(ctx, windowsPath, "-incognito", url)
	}
	return open(ctx, windowsPath, url)
}
