package streamdeck

import "os"

func shutdownSignals() []os.Signal {
	return append([]os.Signal{os.Interrupt}, extraShutdownSignals...)
}
