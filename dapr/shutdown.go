// +build !windows

package dapr

import (
	"os"
	"os/signal"
	"syscall"
)

func SetupShutdownNotify(sigCh chan os.Signal) {
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
}
