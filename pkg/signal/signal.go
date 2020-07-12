package signal

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var once sync.Once //nolint:gochecknoglobals

func Handle(callback func()) {
	once.Do(func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(
			signalChan, syscall.SIGHUP, syscall.SIGINT,
			syscall.SIGTERM, syscall.SIGQUIT)
		<-signalChan
		callback()
	})
}
