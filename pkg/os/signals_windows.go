//go:build windows

package os

import "os"

// On Windows, sending os.Interrupt to a process with os.Process.Signal is not
// implemented; it will return an error instead of sending a signal
var Signals = []os.Signal{os.Interrupt}