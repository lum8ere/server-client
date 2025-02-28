package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForSignalToShutdown() os.Signal {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	res := <-sigs
	return res
}
