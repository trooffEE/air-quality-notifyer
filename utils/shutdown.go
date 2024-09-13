package utils

import (
	"os"
	"os/signal"
)

func DeferShutdown(callback func()) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	go func() {
		select {
		case <-c:
			callback()
			os.Exit(1)
		}
	}()
}
