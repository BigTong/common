package signal

import (
	"os"
	"os/signal"
	"syscall"
)

func Handle(handler func()) {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		for {
			select {
			case <-c:
				handler()
				os.Exit(0)
			}
		}
	}()
}
