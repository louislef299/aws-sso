package browser

import (
	"context"
	"errors"
	"log"
	"strings"

	b "github.com/pkg/browser"
)

var (
	ErrBrowserUnimplemented = errors.New("browser type unimplemented for your OS")
	ErrPrivateNotSupported  = errors.New("private call not supported for browser")
)

type Browser interface {
	OpenURL(ctx context.Context, url string) error
	Type() string
}

func GetBrowser(browserName string, private bool) Browser {
	switch strings.ToLower(browserName) {
	case "chrome":
		return &Chrome{
			private: private,
		}
	case "firefox":
		return &Firefox{
			private: private,
		}
	default:
		return &Default{
			private: private,
		}
	}
}

type Default struct {
	private bool
}

func (d *Default) OpenURL(ctx context.Context, url string) error {
	if d.private {
		log.Println("WARNING: Opening the default browser in private mode is not supported")
	}
	return b.OpenURL(url)
}

func (d *Default) Type() string { return "default" }
