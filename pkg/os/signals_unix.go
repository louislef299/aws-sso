//go:build linux || darwin

package os

import (
	"os"
	"syscall"
)

var Signals = []os.Signal{
	// Interrupt from keyboard
	os.Interrupt,

	// Hangup detected on controlling terminal or death of controlling process
	syscall.SIGHUP,

	// Quit from keyboard
	syscall.SIGQUIT,

	// Termination signal
	syscall.SIGTERM,
}
