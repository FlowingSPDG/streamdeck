//go:build unix

package streamdeck

import (
	"os"
	"syscall"
)

var extraShutdownSignals = []os.Signal{syscall.SIGTERM}
