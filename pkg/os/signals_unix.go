//go:build linux || darwin

package os

import (
	"os"
	"syscall"
)

var Signals = []os.Signal{
	os.Interrupt,
	syscall.SIGHUP,
	syscall.SIGTERM,
	syscall.SIGQUIT,
}
