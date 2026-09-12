//go:build !unix

package streamdeck

import "os"

var extraShutdownSignals []os.Signal
